"Kotlin lint configuration for xquare-application."

load("@aspect_rules_lint//lint:ktlint.bzl", "lint_ktlint_aspect")
load("@aspect_rules_lint//lint:lint_test.bzl", "lint_test")

ktlint = lint_ktlint_aspect(
    binary = Label("@xquare_application_ktlint//file"),
    editorconfig = Label("//systems/xquare-application:.editorconfig"),
    # rules_lint requires a label even when the project has no baselined violations.
    baseline_file = Label("//systems/xquare-application:ktlint-baseline.xml"),
    rule_kinds = [
        "kt_jvm_binary",
        "kt_jvm_library",
        "kt_jvm_test",
    ],
)

ktlint_test = lint_test(aspect = ktlint)
