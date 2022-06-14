resource "kubeberth_server" "terraform-example" {
  name        = "terraform-example"
  running     = true
  cpu         = "1"
  memory      = "1Gi"
  mac_address = "52:42:00:10:00:00"
  hostname    = "terraform-example-server"
  disk        = "terraform-example"
  cloudinit   = "terraform-example"
}
