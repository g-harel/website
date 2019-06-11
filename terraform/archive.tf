data "archive_file" "build_function" {
  type        = "zip"
  source_dir  = "../functions"
  output_path = ".temp/build.zip"
}

data "archive_file" "build_function_templates" {
  type        = "zip"
  source_dir  = "../templates"
  output_path = ".temp/templates.zip"
}
