resource "kubeberth_server" "terraform-example" {
  name        = "terraform-example"
  running     = true
  cpu         = 2
  memory      = "2Gi"
  mac_address = "52:42:00:11:22:33"
  hostname    = "terraform-example-server"
  hosting     = "node-1.k8s.home.arpa"
  disk        = {
    name = "terraformexaample"
  }
  cloudinit   = {
    name = "terraform-example"
  }
}
