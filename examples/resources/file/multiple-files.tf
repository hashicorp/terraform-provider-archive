# Archive multiple files and exclude file.

resource "archive_file" "dotfiles" {
  type        = "zip"
  output_path = "${path.module}/files/dotfiles.zip"
  excludes    = ["${path.module}/unwanted.zip"]

  source {
    content  = data.template_file.vimrc.rendered
    filename = ".vimrc"
  }

  source {
    content  = data.template_file.ssh_config.rendered
    filename = ".ssh/config"
  }
}
