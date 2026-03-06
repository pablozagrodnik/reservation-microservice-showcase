terraform {
	required_providers {
		digitalocean =  {
			source = "digitalocean/digitalocean"
			version = "~> 2.0"
		}
	}

	backend "s3" {
		endpoints = {
			s3 = "https://fra1.digitaloceanspaces.com"
		}
		key							= "terraform.tfstate"
		bucket						= "paws-terraform-state"
		region						= "us-east-1"
		skip_credentials_validation = true
		skip_metadata_api_check 	= true
		skip_requesting_account_id	= true
		skip_s3_checksum			= true
	
	}
}

provider "digitalocean" {
	token = var.do_token
}
