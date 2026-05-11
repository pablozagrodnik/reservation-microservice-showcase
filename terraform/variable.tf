variable "do_token" {
	type		= string
	description	= "DigitalOcean API Token"
}

variable "image_tag" {
	type 		= string
	description = "Docker image tag (Git SHA)"
}

variable "db_password" {
	type      = string
	sensitive = true
}

variable "db_private_host" {
	type = string
}