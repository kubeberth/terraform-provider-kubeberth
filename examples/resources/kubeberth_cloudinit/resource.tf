resource "kubeberth_cloudinit" "terraform-example" {
  name      = "terraform-example"
  user_data = <<EOF
    #cloud-config
    timezone: Asia/Tokyo
    ssh_pwauth: True
    password: ubuntu
    chpasswd: { expire: False }
    disable_root: false
    ssh_authorized_keys:
    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzlOwyoT8qOMpkb9TafGPSM8lXxjZgvIAwHyhNLF1OUOBe8w55KMQ0IR6Q5w1lkKTmMsx7294Fd+xe5ak1BfuwwtF8eOcvWWibDyOr/aPmCFT/N6sZVe2BmN756U1PNDzhufNBH0Yq/AWpZsYn4EQL68hKZuUlA8awOBZS/EfZyPLLCNN5sGSo9nGTBT8DWnC6cEzWJ7ZrBuC69sInYF3haItnYVlafbus07H7waca6WXqZJUpeW0A8Acvsp2EUhNl8Kng/nlnnW4TuuccIGgTNn0Hx1QF6dnLMibD3uqkfAz2QBkJES4K3WWKApGJQxP6h4tw6llDrX7l6m7vHZpn
  EOF
}

