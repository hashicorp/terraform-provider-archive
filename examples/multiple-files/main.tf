# Archive multiple files from a template.

terraform {
  required_providers {
    template = {
      source  = "hashicorp/template"
      version = "2.2.0"
    }
  }
}

data "template_file" "foo" {
  template = file("${path.module}/foo.tpl")
  vars = {
    foo = "bar"
  }
}

data "template_file" "hello" {
  template = file("${path.module}/hello.tpl")
  vars = {
    hello = "world"
  }
}

data "archive_file" "example" {
  type        = "zip"
  output_path = "${path.module}/example.zip"

  source {
    content  = data.template_file.foo.rendered
    filename = "foo.txt"
  }

  source {
    content  = data.template_file.hello.rendered
    filename = "hello.txt"
  }
}
