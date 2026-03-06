import { ChevronRight } from "lucide-react";
import { Link } from "@/i18n/routing";
import { getTranslations } from "next-intl/server";
import { getGeneralDocs } from "@/lib/registry";
import { MDXRemote } from 'next-mdx-remote/rsc';
import { useMDXComponents } from "@/components/mdx-components";
import { TableOfContents } from "@/components/table-of-contents";
import remarkGfm from 'remark-gfm';

export default async function UnderstandingComponentsPage({ params }: { params: Promise<{ locale: string }> }) {
    const resolvedParams = await params;
    const { locale } = resolvedParams;
    const t = await getTranslations({ locale, namespace: "Docs" });

    // Use translations as markdown source
    const markdownSource = `
${t("uc_p1")}

## ${t("uc_how_title")}
${t("uc_how_desc")}

${t("uc_how_desc2")}

## ${t("uc_use_title")}
${t("uc_use_desc")}
- ${t("uc_benefit1")}
- ${t("uc_benefit2")}
- ${t("uc_benefit3")}

## ${t("uc_summary_title")}
${t("uc_summary")}
    `;

    return (
        <main className="relative py-6 lg:py-8 w-full max-w-3xl mx-auto md:mx-0">
            <div className="mx-auto w-full min-w-0 pb-16">
                <div className="mb-4 flex items-center space-x-1 text-sm text-muted-foreground">
                    <Link className="overflow-hidden text-ellipsis whitespace-nowrap hover:underline" href="/docs">Docs</Link>
                    <ChevronRight className="h-4 w-4" />
                    <div className="font-medium text-foreground">{t("uc_title")}</div>
                </div>

                <div className="space-y-2">
                    <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-10">{t("uc_title")}</h1>
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
