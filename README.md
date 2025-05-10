# 🚀 DockDev: Your Docker Development Companion

![DockDev](https://img.shields.io/badge/DockDev-CLI%20Tool-brightgreen.svg)
![Docker](https://img.shields.io/badge/Docker-%F0%9F%90%B3-blue.svg)
![NGINX](https://img.shields.io/badge/NGINX-%F0%9F%93%8E-orange.svg)

Welcome to **DockDev**, a powerful command-line interface tool designed to create isolated Docker-based development environments. With DockDev, you can easily set up local domain access using NGINX as a reverse proxy. This tool is specifically tailored for Windows users utilizing WSL2.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Releases](#releases)

## Features

- **Isolated Development Environments**: Quickly spin up separate environments for different projects.
- **Local Domain Access**: Access your applications via local domains, enhancing your development experience.
- **NGINX Reverse Proxy**: Utilize NGINX to manage your local traffic efficiently.
- **Windows + WSL2 Compatibility**: Seamlessly integrate with Windows Subsystem for Linux 2.
- **Simple CLI**: Use a straightforward command-line interface for all operations.

## Installation

To get started with DockDev, follow these simple steps:

1. **Download the latest release** from the [Releases section](https://github.com/PRATHISHX2004/dockdev/releases). You will need to download and execute the relevant file for your system.
2. **Extract the files** to your preferred directory.
3. **Add DockDev to your PATH** to run it from anywhere in your terminal.

### Prerequisites

- Docker installed on your machine.
- WSL2 set up on Windows.
- Basic knowledge of command-line operations.

## Usage

Once you have installed DockDev, you can start using it right away. Here are some common commands:

### Create a New Environment

To create a new development environment, use the following command:

```bash
dockdev create <environment-name>
```

### Start an Environment

To start your newly created environment, run:

```bash
dockdev start <environment-name>
```

### Stop an Environment

To stop an environment, use:

```bash
dockdev stop <environment-name>
```

### Remove an Environment

If you want to remove an environment, simply run:

```bash
dockdev remove <environment-name>
```

## Configuration

DockDev uses a configuration file to manage your environments. You can find this file in your home directory under `.dockdev/config.yaml`. Here you can set various options such as:

- Default ports
- NGINX configurations
- Local domain names

### Example Configuration

```yaml
default_ports:
  - 80
  - 443
nginx:
  server_name: "local.dev"
  root: "/var/www/html"
```

## Examples

### Setting Up a Basic Environment

1. Create a new environment:

   ```bash
   dockdev create myapp
   ```

2. Start the environment:

   ```bash
   dockdev start myapp
   ```

3. Access your application at `http://local.dev`.

### Using Multiple Environments

You can create multiple environments for different projects. Just ensure each has a unique name and domain.

```bash
dockdev create project1
dockdev create project2
```

Start them individually as needed.

## Contributing

We welcome contributions! If you want to help improve DockDev, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your branch and submit a pull request.

## License

DockDev is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any inquiries, please reach out via GitHub issues or directly at [your-email@example.com](mailto:your-email@example.com).

## Releases

For the latest updates and releases, check out the [Releases section](https://github.com/PRATHISHX2004/dockdev/releases). You will need to download and execute the relevant file for your system.

---

Thank you for using DockDev! Happy coding!