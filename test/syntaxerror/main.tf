resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "sleep ${var.sleep_time}"
    interpreter = ["/bin/sh", "-c"]
  }
# Intentional syntax error
