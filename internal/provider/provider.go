// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/keeper-security/keeper-sdk-golang/sdk/api"
	"github.com/keeper-security/keeper-sdk-golang/sdk/auth"
	"github.com/keeper-security/keeper-sdk-golang/sdk/enterprise"
	"github.com/keeper-security/keeper-sdk-golang/sdk/helpers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

// Ensure KeeperEnterpriseProvider satisfies various provider interfaces.
var _ provider.Provider = &keeperEnterpriseProvider{}

// KeeperEnterpriseProvider defines the provider implementation.
type keeperEnterpriseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &keeperEnterpriseProvider{
			version: version,
		}
	}
}

type keeperEnterpriseProviderModel struct {
	ConfigurationPath types.String `tfsdk:"config_path"`
	ConfigurationType types.String `tfsdk:"config_type"`
}

func (p *keeperEnterpriseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kepr"
	resp.Version = p.version
}

func (p *keeperEnterpriseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config_path": schema.StringAttribute{
				MarkdownDescription: "Full path to the Keeper configuration path",
				Optional:            true,
			},
			"config_type": schema.StringAttribute{
				MarkdownDescription: "Configuration file type: sdk | commander",
				Optional:            true,
			},
		},
	}
}
func (p *keeperEnterpriseProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var config keeperEnterpriseProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ConfigurationType.IsNull() || config.ConfigurationType.IsUnknown() {
		return
	}
	var configType = config.ConfigurationType.ValueString()
	if strings.EqualFold(configType, "sdk") || strings.EqualFold(configType, "commander") {
		return
	}
	resp.Diagnostics.AddAttributeError(
		path.Root("config_type"),
		fmt.Sprintf("Invalid Configuration Attribute Value (%s)", configType),
		"Expected either \"sdk\" or \"commander\"",
	)
}

func (p *keeperEnterpriseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config keeperEnterpriseProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ConfigurationPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_path"),
			"Unknown Attribute Value",
			"The provider uses configuration file to connect to the Keeper backend.",
		)
	}
	if config.ConfigurationType.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_type"),
			"Unknown Attribute Value",
			"The provider uses configuration file to connect to the Keeper backend.",
		)
	}

	var core = zapcore.RegisterHooks(NewNopCore(), func(entry zapcore.Entry) error {
		switch entry.Level {
		case zapcore.DebugLevel:
			tflog.Debug(ctx, entry.Message)
		case zapcore.InfoLevel:
			tflog.Info(ctx, entry.Message)
		case zapcore.WarnLevel:
			tflog.Warn(ctx, entry.Message)
		case zapcore.ErrorLevel:
			tflog.Error(ctx, entry.Message)
		}
		return nil
	})

	var logger = zap.New(core)
	api.SetLogger(logger)
	var configFilename string
	if config.ConfigurationPath.IsNull() {
		configFilename = "config.json"
	} else {
		configFilename = config.ConfigurationPath.ValueString()
	}
	configFilename = api.GetKeeperFileFullPath(configFilename)
	tflog.Info(ctx, "Configuring file path: "+configFilename)
	var err error
	if _, err = os.Stat(configFilename); err != nil {
		return
	}
	var isCommanderConfig = false
	if !config.ConfigurationType.IsNull() {
		var configType = config.ConfigurationType.ValueString()
		isCommanderConfig = strings.EqualFold(configType, "commander")
	}
	var configStorage auth.IConfigurationStorage
	if isCommanderConfig {
		configStorage = helpers.NewCommanderConfiguration(configFilename)
	} else {
		configStorage = helpers.NewJsonConfigurationFile(configFilename)
	}
	var keeperConfig auth.IKeeperConfiguration
	if keeperConfig, err = configStorage.Get(); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_path"),
			"Keeper configuration file does not exist",
			"The provider requires configuration file to connect to the Keeper backend.",
		)
	}
	var keeperEndpoint = helpers.NewKeeperEndpoint(keeperConfig.LastServer(), configStorage)
	var loginAuth = helpers.NewLoginAuth(keeperEndpoint)
	loginAuth.Login(keeperConfig.LastLogin())
	var step = loginAuth.Step()
	if step.LoginState() != auth.LoginState_Connected {
		var stepInfo string
		switch step.LoginState() {
		case auth.LoginState_TwoFactor:
			stepInfo = "Requires 2FA"
		case auth.LoginState_DeviceApproval:
			stepInfo = "Requires Device Approval"
		case auth.LoginState_Password:
			stepInfo = "Requires Password"
		case auth.LoginState_Error:
			if es, ok := step.(auth.IErrorStep); ok {
				stepInfo = fmt.Sprintf("Error: %s", es.Error().Error())
			} else {
				stepInfo = "Error"
			}
		default:
			stepInfo = "SSO is not supported"
		}
		resp.Diagnostics.AddError(
			"Cannot connect to Keeper in unattended mode",
			"The provider requires configuration file to be prepared for unattended login mode.\n"+
				"It requires either Persistent Login or storing the user password into the configuration file.\n"+
				"Login Step: "+stepInfo,
		)
		return
	}
	var ok bool
	var connectedStep auth.IConnectedStep
	if connectedStep, ok = step.(auth.IConnectedStep); !ok {
		resp.Diagnostics.AddError(
			"Cannot connect to Keeper in unattended mode",
			"Keeper SDK library error.",
		)
		return
	}
	var keeperAuth auth.IKeeperAuth
	if keeperAuth, err = connectedStep.TakeKeeperAuth(); err != nil {
		resp.Diagnostics.AddError(
			"Cannot connect to Keeper in unattended mode",
			err.Error(),
		)
		return
	}
	var loader = enterprise.NewEnterpriseLoader(keeperAuth, nil)
	if err = loader.Load(); err != nil {
		resp.Diagnostics.AddError(
			"Cannot load Keeper enterprise information",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = loader
	resp.ResourceData = loader
}

func (p *keeperEnterpriseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTeamResource,
	}
}

func (p *keeperEnterpriseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNodeDataSource, NewNodesDataSource,
		NewTeamDataSource, NewTeamsDataSource,
		NewUsersDataSource,
	}
}
