resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "sleep ${var.sleep_time}; false"
    interpreter = ["/bin/sh", "-c"]
  }
}

