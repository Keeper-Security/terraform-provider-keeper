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
	"github.com/keeper-security/keeper-sdk-golang/api"
	"github.com/keeper-security/keeper-sdk-golang/auth"
	"github.com/keeper-security/keeper-sdk-golang/enterprise"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"terraform-provider-keeper/internal/acc_test"
	"terraform-provider-keeper/internal/model"
)

// Ensure KeeperEnterpriseProvider satisfies various provider interfaces.
var _ provider.Provider = &keeperEnterpriseProvider{}

// KeeperEnterpriseProvider defines the provider implementation.
type keeperEnterpriseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version    string
	management *enterprise.IEnterpriseManagement
}

func New(version string, management *enterprise.IEnterpriseManagement) func() provider.Provider {
	return func() provider.Provider {
		return &keeperEnterpriseProvider{
			version:    version,
			management: management,
		}
	}
}

type keeperEnterpriseProviderModel struct {
	ConfigurationPath types.String `tfsdk:"config_path"`
	ConfigurationType types.String `tfsdk:"config_type"`
	Password          types.String `tfsdk:"password"`
}

func (p *keeperEnterpriseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "keeper"
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
			"password": schema.StringAttribute{
				MarkdownDescription: "Keeper account password",
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

	tflog.Warn(ctx, fmt.Sprintf("Called Provider configure: %t", *p.management == nil))
	if *p.management == nil {
		var core = model.NewTerraformCore(ctx, zapcore.DebugLevel)
		var logger = zap.New(core)
		api.SetLogger(logger)

		if p.version == "test" {
			*p.management = acc_test.NewTestingManagement()
		} else {
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
				configStorage = auth.NewCommanderConfiguration(configFilename)
			} else {
				configStorage = auth.NewJsonConfigurationFile(configFilename)
			}
			var keeperConfig auth.IKeeperConfiguration
			if keeperConfig, err = configStorage.Get(); err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("config_path"),
					"Keeper configuration file does not exist",
					"The provider requires configuration file to connect to the Keeper backend.",
				)
			}
			var keeperEndpoint = auth.NewKeeperEndpoint(keeperConfig.LastServer(), configStorage)
			var loginAuth = auth.NewLoginAuth(keeperEndpoint)
			var passwords []string
			if !config.Password.IsNull() {
				var passwd = config.Password.ValueString()
				if len(passwd) > 0 {
					passwords = append(passwords, passwd)
				}
			}
			loginAuth.Login(keeperConfig.LastLogin(), passwords...)
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

			*p.management = enterprise.NewSyncEnterpriseManagement(loader)
		}
	}
	resp.DataSourceData = (*p.management).EnterpriseData()
	resp.ResourceData = *p.management
}

func (p *keeperEnterpriseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newTeamResource, newNodeResource, newRoleResource,
		newTeamMembershipResource, newRoleMembershipResource,
		newManagedNodeResource, newRoleEnforcementsResource,
	}
}

func (p *keeperEnterpriseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newNodeDataSource, newNodesDataSource,
		newTeamDataSource, newTeamsDataSource,
		newRoleDataSource, newRolesDataSource,
		newPrivilegeDataSource,
		newEnforcementsDataSource, newEnforcementsAccountDataSource, newEnforcementsLoginDataSource,
		newEnforcementsAllowIpListDataSource, newEnforcementsPlatformDataSource, newEnforcementsSharingDataSource,
		newEnforcements2faDataSource, newEnforcementsKeeperFillDataSource, newEnforcementsRecordTypesDataSource,
		newEnforcementsVaultDataSource,
		newUserDataSource, newUsersDataSource,
	}
}
