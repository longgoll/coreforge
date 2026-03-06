"use client";

import { Link } from "@/i18n/routing";
import { usePathname } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";

const NavItem = ({ href, children, badge, pathname }: { href: string; children: React.ReactNode; badge?: string; pathname: string }) => {
    const isActive = pathname === href;
    return (
        <Link
            href={href}
            className={`group flex w-full items-center justify-between rounded-md border border-transparent px-2 py-1.5 hover:underline text-sm font-medium ${isActive ? 'text-foreground bg-muted/60' : 'text-muted-foreground'}`}
        >
            {children}
            {badge && <Badge variant="secondary" className="h-4 text-[9px] px-1 py-0">{badge}</Badge>}
        </Link>
    );
};

export interface DocsSidebarProps {
    categories: { [category: string]: { id: string, name: string }[] };
}

export function DocsSidebar({ categories }: DocsSidebarProps) {
    const pathname = usePathname();
    const t = useTranslations("DocsSidebar");

    return (
        <aside className="fixed top-14 z-30 -ml-2 hidden h-[calc(100vh-3.5rem)] w-full shrink-0 md:sticky md:block">
            <ScrollArea className="h-full py-6 pr-6 lg:py-8">
                <div className="w-full">
                    <div className="pb-6">
                        <h4 className="mb-2 rounded-md px-2 py-1 text-sm font-semibold capitalize">{t("gettingStarted")}</h4>
                        <div className="grid grid-flow-row auto-rows-max text-sm gap-0.5">
                            <NavItem href="/docs" pathname={pathname}>{t("introduction")}</NavItem>
                            <NavItem href="/docs/understanding-components" pathname={pathname}>{t("understandingComponents")}</NavItem>
                            <NavItem href="/docs/installation" pathname={pathname}>{t("installation")}</NavItem>
                        </div>
                    </div>

                    {Object.keys(categories).sort().map(category => (
                        <div key={category} className="pb-6">
                            <h4 className="mb-2 rounded-md px-2 py-1 text-sm font-semibold capitalize">
                                {t.has(`categories.${category}`) ? t(`categories.${category}`) : category}
                            </h4>
                            <div className="grid grid-flow-row auto-rows-max text-sm gap-0.5">
                                {categories[category].map(comp => (
                                    <NavItem key={comp.id} href={`/docs/components/${comp.id}`} pathname={pathname}>
                                        {t.has(`components.${comp.id}`) ? t(`components.${comp.id}`) : comp.name}
                                    </NavItem>
                                ))}
                            </div>
                        </div>
                    ))}

                </div>
            </ScrollArea>
        </aside>
    );
}
