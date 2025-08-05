provider "keeper" {
  config_path = "c:/tf/config.json"
  # config_path = "~/tf/config.json"
  config_type = "commander" # Default
  password    = "test123"   # Optional if configuration users persistent login
}
