package builder

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
)

//go:embed templates/archlinux.sh
var archlinuxTemplate string

// TargetTypeArchlinux identifies the Archlinux target.
const TargetTypeArchlinux Type = "arch"

func init() {
	BuilderByTarget[TargetTypeArchlinux] = &archlinux{}
}

// archlinux is a driverkit target.
type archlinux struct {
}

type archlinuxTemplateData struct {
	commonTemplateData
	KernelDownloadURL string
}

func (c *archlinux) Name() string {
	return TargetTypeArchlinux.String()
}

func (c *archlinux) TemplateScript() string {
	return archlinuxTemplate
}

func (c *archlinux) URLs(kr kernelrelease.KernelRelease) ([]string, error) {

	urls := []string{}
	possibleCompressionSuffixes := []string{
		"xz",
		"zst",
	}

	// check the architecture, which limits the mirror options
	if kr.Architecture.ToNonDeb() == "x86_64" {
		if strings.Contains(kr.FullExtraversion, "arch") { // arch stable kernel
			baseURL := "https://archive.archlinux.org/packages/l/linux-headers"
			for _, compressionAlgo := range possibleCompressionSuffixes {
				urls = append(
					urls,
					fmt.Sprintf(
						"%s/linux-headers-%s%s-%s.pkg.tar.%s",
						baseURL,
						kr.Fullversion,
						kr.FullExtraversion,
						kr.Architecture.ToNonDeb(),
						compressionAlgo,
					),
				)
			}
		} else if strings.Contains(kr.FullExtraversion, "hardened") || strings.Contains(kr.FullExtraversion, ".a-1") { // arch hardened kernel ("a-1" is old naming standard)
			baseURL := "https://archive.archlinux.org/packages/l/linux-hardened-headers"
			for _, compressionAlgo := range possibleCompressionSuffixes {
				urls = append(
					urls,
					fmt.Sprintf(
						"%s/linux-hardened-headers-%s%s-%s.pkg.tar.%s",
						baseURL,
						kr.Fullversion,
						kr.FullExtraversion,
						kr.Architecture.ToNonDeb(),
						compressionAlgo,
					),
				)
			}
		} else if strings.Contains(kr.FullExtraversion, "zen") { // arch zen kernel
			baseURL := "https://archive.archlinux.org/packages/l/linux-zen-headers"
			for _, compressionAlgo := range possibleCompressionSuffixes {
				urls = append(
					urls,
					fmt.Sprintf(
						"%s/linux-zen-headers-%s%s-%s.pkg.tar.%s",
						baseURL,
						kr.Fullversion,
						kr.FullExtraversion,
						kr.Architecture.ToNonDeb(),
						compressionAlgo,
					),
				)
			}
		} else { // arch LTS kernel
			baseURL := "https://archive.archlinux.org/packages/l/linux-lts-headers"
			for _, compressionAlgo := range possibleCompressionSuffixes {
				urls = append(
					urls,
					fmt.Sprintf(
						"%s/linux-lts-headers-%s%s-%s.pkg.tar.%s",
						baseURL,
						kr.Fullversion,
						kr.FullExtraversion,
						kr.Architecture.ToNonDeb(),
						compressionAlgo,
					),
				)
			}
		}
	} else if kr.Architecture.ToNonDeb() == "aarch64" {
		baseURL := "http://tardis.tiny-vps.com/aarm/packages/l/linux-aarch64-headers/"
		for _, compressionAlgo := range possibleCompressionSuffixes {
			urls = append(
				urls,
				fmt.Sprintf(
					"%s/linux-aarch64-headers-%s%s-%s.pkg.tar.%s",
					baseURL,
					kr.Fullversion,
					kr.FullExtraversion,
					kr.Architecture.ToNonDeb(),
					compressionAlgo,
				),
			)
		}
	}

	return urls, nil
}

func (c *archlinux) TemplateData(cfg Config, kr kernelrelease.KernelRelease, urls []string) interface{} {
	return archlinuxTemplateData{
		commonTemplateData: cfg.toTemplateData(c, kr),
		KernelDownloadURL:  urls[0],
	}
}
