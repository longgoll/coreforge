import { DocsSidebar } from "@/components/docs-sidebar";
import { getRegistry } from "@/lib/registry";

export default async function DocsLayout({ children }: { children: React.ReactNode }) {
    const registry = await getRegistry();
    const categories: { [key: string]: { id: string, name: string }[] } = {};

    for (const [id, comp] of Object.entries(registry.components)) {
        const catName = comp.category || "General";
        if (!categories[catName]) {
            categories[catName] = [];
        }
        categories[catName].push({
            id,
            name: id.split('-').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
        });
    }

    return (
        <div className="container mx-auto flex-1 items-start md:grid md:grid-cols-[220px_minmax(0,1fr)] lg:grid-cols-[240px_minmax(0,1fr)] xl:grid-cols-[240px_minmax(0,1fr)_250px] gap-6 lg:gap-10 px-4 sm:px-8">
            <DocsSidebar categories={categories} />
            {children}
        </div>
    );
}
