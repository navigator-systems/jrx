# JRX - Project Management CLI

![JRX Logo](images/logo.png)

**JRX** is a simple command-line tool for project management and template-based project generation. It provides a streamlined way to create new projects from predefined templates and manage project information.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/navigator-systems/jrx.git
cd jrx

# Build and install
make install
```

### Build Options

```bash
# Development build
make dev

# Production build (optimized)
make prod

# Cross-platform compilation
make compile
```

## Features

- **Template-based Project Generation**: Create new projects from predefined templates
- **Project Information**: Get detailed information about existing projects  
- **Template Management**: Download and manage project templates
- **Git Integration**: Automatic Git repository initialization for new projects

## Usage

JRX provides two main command groups: `project` and `templates`.

### Project Commands

#### Create a New Project

```bash
# Create a new project from a template
jrx project new <project-name> <template-name>

# Example: Create a Go web service
jrx project new my-web-app golang-web

```

### Template Commands

#### Download Templates

```bash
# Download/update template repository
jrx templates download
```

#### List Available Templates

```bash
# List all available templates
jrx templates list

# Aliases
jrx t list
```

## Template Features

All templates support:
- **Variable Substitution**: Dynamic content based on project configuration
- **Directory Structure**: Maintains proper project organization
- **Git Integration**: Automatic repository initialization with main branch
- **Configuration Files**: Project-specific settings and metadata

## Configuration

JRX uses TOML configuration files for template management and project settings.

### Global Configuration (.jrxrc)

JRX reads a global configuration file located at `~/.jrxrc` (in your home directory) that controls where templates are hosted and how to access them. This file should contain:

```toml
templates_repo = "git@github.com:your-org/jrx-templates.git"
templates_branch = "main"
ssh_key_path = "/home/user/.ssh/id_rsa"
ssh_key_passphrase = "your-passphrase" 
```

**Configuration Options:**
- `templates_repo`: Git repository URL containing your custom templates
- `templates_branch`: Branch to use from the templates repository (e.g., "main", "develop")
- `ssh_key_path`: Path to your SSH private key for accessing private repositories
- `ssh_key_passphrase`: Passphrase for your SSH key (optional if key has no passphrase)


### Template Configuration

The main template configuration is located in `jrxTemplates/templates.toml` after downloading templates.


## Dependencies

- Go 1.25+
- Git (for project initialization and template management)

## Build Requirements

- Go toolchain 1.25
- Make (for build automation)

## Command Reference

```bash
NAME:
   jrx - Just a simple project management CLI

USAGE:
   jrx [global options] command [command options]

COMMANDS:
   project, p    Manage projects
     new, n      Create a new project
   templates, t  Manage project templates
     list        Get information about templates
     download    Download the templates for a new project
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the terms specified in the LICENSE file.