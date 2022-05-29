resource "kubeberth_disk" "terraform-example" {
  name    = "terraform-example"
  size    = "16Gi"
  archive = "terraform-example"
}
