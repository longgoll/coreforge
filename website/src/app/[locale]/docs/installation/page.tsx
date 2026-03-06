import { ChevronRight } from "lucide-react";
import { Link } from "@/i18n/routing";
import { getTranslations } from "next-intl/server";
import { getGeneralDocs } from "@/lib/registry";
import { MDXRemote } from 'next-mdx-remote/rsc';
import { useMDXComponents } from "@/components/mdx-components";
import { TableOfContents } from "@/components/table-of-contents";
import remarkGfm from 'remark-gfm';

export default async function InstallationPage({ params }: { params: Promise<{ locale: string }> }) {
    const resolvedParams = await params;
    const { locale } = resolvedParams;
    const t = await getTranslations({ locale, namespace: "Docs" });

    // Use translations as markdown source
    const markdownSource = `
${t("install_p1")}

## ${t("install_req")}
${t("install_req_desc")}
- Go 1.22+
- Node.js 18+ / Go 1.22+ / .NET 8.0+

## ${t("install_step1")}
${t("install_step1_desc")}

\`\`\`bash
go install github.com/longgoll/coreforge/cmd/forge@latest
\`\`\`

## ${t("install_step2")}
${t("install_step2_desc")}

\`\`\`bash
forge init
\`\`\`

## ${t("install_step3")}
${t("install_step3_desc")}

\`\`\`bash
forge add error-handler
\`\`\`

## ${t("install_step4")}
${t("install_step4_desc")}
    `;

    return (
        <main className="relative py-6 lg:py-8 w-full max-w-3xl mx-auto md:mx-0">
            <div className="mx-auto w-full min-w-0 pb-16">
                <div className="mb-4 flex items-center space-x-1 text-sm text-muted-foreground">
                    <Link className="overflow-hidden text-ellipsis whitespace-nowrap hover:underline" href="/docs">Docs</Link>
                    <ChevronRight className="h-4 w-4" />
                    <div className="font-medium text-foreground">{t("install_title")}</div>
                </div>

                <div className="space-y-2">
                    <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-10">{t("install_title")}</h1>
                </div>

                <div className="mb-14">
                    <MDXRemote
                        source={markdownSource}
                        components={useMDXComponents({})}
                        options={{ mdxOptions: { remarkPlugins: [remarkGfm] } }}
                    />
                </div>
            </div>

            <div className="hidden xl:block">
                <div className="sticky top-20 h-[calc(100vh-5rem)] overflow-y-auto pt-8">
                    <TableOfContents />
                </div>
            </div>
        </main>
    );
}
