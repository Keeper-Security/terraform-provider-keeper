data "keeper_enforcements_record_types" "example" {
  restrict_record_types = [
    "address", "bankAccount", "bankCard", "birthCertificate", "contact", "databaseCredentials",
    "driverLicense", "encryptedNotes", "file", "general", "healthInsurance", "login", "membership",
    "passport", "photo", "serverCredentials", "softwareLicense", "sshKeys", "ssnCard"
  ]
}