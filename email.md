Hello {{ .Customer }}!

Prior to our next meeting, please download and transfer this software to the correct server. You will notice repeated download links in the list below, the actual software needs to be downloaded once, but I included all the relevant links to make it clear where you need to transfer the software. 

- Package Manager Required Software
  - Package Manager: {{ .PackageManager }} {{ range .R}}{{.}}{{end}}
- Connect Required Software
  - Connect: {{ .Connect }} {{ range .R}}{{.}}{{end}} {{ range .Python}}{{.}}{{end}}
  - Quarto:
  - Professional Drivers: {{ .ProDriver }}
- Workbench Required Software
  - Workbench: {{ .Workbench }} {{ range .R}}{{.}}{{end}} {{ range .Python}}{{.}}{{end}}
  - Quarto:
  - Professional Drivers: {{ .ProDriver }}

When we start meeting, we will install the products starting with Package Manager and make our way down the list. This is because Package Manager needs to be setup in order to correctly configure both Connect and Workbench. 

Thanks!

