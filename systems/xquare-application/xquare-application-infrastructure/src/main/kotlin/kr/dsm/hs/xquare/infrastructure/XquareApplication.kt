package kr.dsm.hs.xquare.infrastructure

import org.springframework.boot.autoconfigure.SpringBootApplication

@SpringBootApplication(
    proxyBeanMethods = false,
    scanBasePackages = ["kr.dsm.hs.xquare"],
)
class XquareApplication
