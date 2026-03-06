import { Link } from "@/i18n/routing";
import { useTranslations } from "next-intl";

export function SiteFooter() {
    const t = useTranslations("Footer");
    return (
        <footer className="py-6 md:px-8 md:py-0 border-t bg-muted/20">
            <div className="container mx-auto px-4 flex flex-col items-center justify-between gap-4 md:h-16 md:flex-row">
                <p className="text-balance text-center text-sm leading-loose text-muted-foreground md:text-left">
                    {t("builtWith")} <Link href="#" className="font-medium underline underline-offset-4 text-foreground">CoreForge</Link>.
                </p>
            </div>
        </footer>
    );
}
