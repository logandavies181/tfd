resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "sleep ${var.sleep_time}"
    interpreter = ["/bin/sh", "-c"]
  }
}

output "sleep_time" {
  value = var.sleep_time
}
