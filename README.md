# Getting started

This repository has the codebase that can parse templates and make replacements by reading a config file (template.yaml). Steps below showcase closing a sample repository (similar to a component repository), ready the template.yaml file and making necessary replacements across package.json, README etc. files found in the repository

```bash
# Run by executing the following command, replace DIR_NAME_OF_YOUR_CHOICE by a directory name such as 'sample', 'test' or others as you like

go run main.go [DIR_NAME_OF_YOUR_CHOICE] https://github.com/pyljain/sample_template.git 
