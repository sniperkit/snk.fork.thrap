manifest "thrap" {
  name = "thrap"

  components {
    nomad {
      // Image name
      name = "nomad"

      // Image version
      version = "0.8.4"

      // Type of component used for automation
      type = "api"

      // Ports the components listens on
      ports {
        http     = 4646
        port4647 = 4647
        port4648 = 4648
      }

      // Env. vars. needed by the component
      env {
        vars {
          CONSUL_ADDR      = "http://${comp.consul.container.addr.http}"
          VAULT_ADDR       = "http://${comp.vault.container.addr.default}"
          BOOTSTRAP_EXPECT = "1"
        }
      }

      // Data needed to build the component
      build {
        dockerfile = "nomad.dockerfile"
      }
    }

    consul {
      name    = "consul"
      version = "1.2.0"
      type    = "api"

      ports {
        http = 8500
      }
    }

    vault {
      name    = "vault"
      version = "0.10.3"
      type    = "api"

      ports {
        default = 8200
      }

      env {
        vars {
          VAULT_DEV_ROOT_TOKEN_ID = "myroot"
        }
      }
    }

    registry {
      // Image name that will be used. 
      // The final full image name will be <stack.id>/<name>:<stack.version>
      name = "registry"

      type     = "api"
      language = "go"

      build {
        dockerfile = "api.dockerfile"
      }

      cmd  = "thrap"
      args = ["agent"]

      ports {
        http = 10000
      }

      secrets {
        // Destination of secrets relative to working dir.
        destination = ".thrap/creds.hcl"

        // Secrets inline template
        template = <<EOF
registry {
  ecr {
    key    = "${aws_access_key_id}"
    secret = "${aws_secret_access_key}"
  }
}

vcs {
  github {
    token = "${github_token}"
  }
}
EOF
      }

      // Head of the stack i.e. api or ui or application exposed interface
      head = true

      env {
        file = ".env"

        vars {
          # Should be available by default  
          STACK_VERSION = "${stack.version}"
          VAULT_ADDR    = "http://${comp.vault.container.addr.default}"
          NOMAD_ADDR    = "http://${comp.nomad.container.addr.http}"
        }
      }
    }
  }

  dependencies {
    github {
      name     = "github"
      version  = "v3"
      external = true
      config   = {}
    }

    ecr {
      external = true
    }

    vault {
      name    = "vault"
      version = "0.10.3"
    }

    consul {
      name    = "consul"
      version = "1.2.0"
    }

    nomad {
      name    = "nomad"
      version = "0.8.4"
    }

    docker {
      name    = "docker"
      version = "1.37"
    }
  }
}
