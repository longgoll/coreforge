import { Link } from '@/i18n/routing';
import { Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { getRegistry, getComponentDocs } from '@/lib/registry';
import { getTranslations } from 'next-intl/server';

function parseFrontmatter(markdown: string) {
    const frontmatterRegex = /^---\s*\n([\s\S]*?)\n---/;
    const match = markdown.match(frontmatterRegex);
    if (!match) return null;

    const data: Record<string, string> = {};
    const lines = match[1].split('\n');
    for (const line of lines) {
        const colonIndex = line.indexOf(':');
        if (colonIndex !== -1) {
            const key = line.slice(0, colonIndex).trim();
            const value = line.slice(colonIndex + 1).trim();
            data[key] = value;
        }
    }
    return data;
}

export default async function ComponentsPage({ params }: { params: Promise<{ locale: string }> }) {
    const resolvedParams = await params;
    const locale = resolvedParams.locale;
    const t = await getTranslations('ComponentsList');

    const registry = await getRegistry();

    // Resolve components asynchronously to allow fetching markdown
    const componentsList = await Promise.all(Object.entries(registry.components).map(async ([id, comp]) => {
        let name = id.split('-').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
        let description = comp.description;

        // Try to fetch localized markdown
        const md = await getComponentDocs(id, locale);
        if (md) {
            const fm = parseFrontmatter(md);
            if (fm) {
                if (fm.title) name = fm.title;
                if (fm.description) description = fm.description;
            }
        }

        return {
            id,
            name,
            description,
            category: comp.category || 'general'
        };
    }));

    // Group components by Category
    const groupedComponents = componentsList.reduce((acc, comp) => {
        const cat = comp.category;
        if (!acc[cat]) {
            acc[cat] = [];
        }
        acc[cat].push(comp);
        return acc;
    }, {} as Record<string, typeof componentsList>);

    // Ensure we have a predictable order (e.g., Auth, Middleware, Utility...)
    const categoryOrder = ["auth", "middleware", "utility", "database", "documentation", "general"];
    const sortedCategories = Object.keys(groupedComponents).sort((a, b) => {
        const indexA = categoryOrder.indexOf(a.toLowerCase());
        const indexB = categoryOrder.indexOf(b.toLowerCase());
        if (indexA !== -1 && indexB !== -1) return indexA - indexB;
        if (indexA !== -1) return -1;
        if (indexB !== -1) return 1;
        return a.localeCompare(b);
    });

    return (
        <div className="container mx-auto py-12 px-4 sm:px-8 max-w-[1200px]">
            <div className="flex flex-col items-start gap-4 md:flex-row md:justify-between md:gap-8 mb-16 text-left">
                <div className="flex-1 space-y-4">
                    <h1 className="inline-block font-heading text-4xl tracking-tight lg:text-5xl font-extrabold">{t('title')}</h1>
                    <p className="text-xl text-muted-foreground leading-relaxed">{t('description')}</p>
                </div>
            </div>

            <div className="flex flex-col gap-16">
                {sortedCategories.map(category => (
                    <section key={category} className="flex flex-col gap-6">
                        <div className="flex items-center gap-4 border-b border-border/50 pb-3">
                            <h2 className="text-2xl font-bold tracking-tight text-foreground/90 capitalize">
                                {t(`categories.${category.toLowerCase()}`) || category}
                            </h2>
                            <Badge variant="outline" className="rounded-full bg-muted/40 font-mono text-xs">
                                {groupedComponents[category].length}
                            </Badge>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                            {groupedComponents[category].map(comp => (
                                <Link href={`/docs/components/${comp.id}`} key={comp.id}>
                                    <Card className="h-full bg-card hover:bg-muted/30 border-border/50 hover:border-primary/40 transition-all duration-300 group cursor-pointer shadow-sm hover:shadow-md overflow-hidden relative">
                                        <div className="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-transparent via-primary/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                                        <CardHeader className="p-5">
                                            <div className="flex justify-between items-start mb-3">
                                                <CardTitle className="text-lg font-semibold text-foreground/80 group-hover:text-primary transition-colors">{comp.name}</CardTitle>
                                            </div>
                                            <CardDescription className="text-sm leading-relaxed text-muted-foreground line-clamp-2">{comp.description}</CardDescription>
                                        </CardHeader>
                                    </Card>
                                </Link>
                            ))}
                        </div>
                    </section>
                ))}
            </div>
        </div>
    );
}
