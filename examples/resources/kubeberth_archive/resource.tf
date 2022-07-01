resource "kubeberth_archive" "terraform-example" {
  name       = "terraform-example"
  repository = "http://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img"
}
