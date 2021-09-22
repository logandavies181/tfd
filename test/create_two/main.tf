resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "sleep ${var.sleep_time}"
    interpreter = ["/bin/sh", "-c"]
  }
}

resource "null_resource" "test2" {
  provisioner "local-exec" {
    command = "sleep ${var.sleep_time}"
    interpreter = ["/bin/sh", "-c"]
  }
}
