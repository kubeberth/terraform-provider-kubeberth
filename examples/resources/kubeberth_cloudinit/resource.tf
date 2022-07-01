resource "kubeberth_cloudinit" "terraform-example" {
  name      = "terraform-example"
  network_data = ""
  user_data = <<EOF
#cloud-config
timezone: Asia/Tokyo
ssh_pwauth: True
password: ubuntu
chpasswd: { expire: False }
disable_root: false
#ssh_authorized_keys:
#- ssh-rsa XXXXXXXXXXXXXXXXXXXXXXXXX
EOF
}
