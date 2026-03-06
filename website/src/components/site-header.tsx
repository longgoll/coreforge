"use client";

import { Link } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { Package2, Search, Github, Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { LanguageSwitcher } from "@/components/language-switcher";

export function SiteHeader() {
    const { setTheme, theme } = useTheme();
    const t = useTranslations("Navigation");

    return (
        <header className="sticky top-0 z-50 w-full border-b bg-background/80 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="container mx-auto flex h-14 items-center px-4 sm:px-8">
                <div className="mr-4 hidden md:flex">
                    <Link href="/" className="mr-6 flex items-center space-x-2">
                        <Package2 className="h-6 w-6 text-primary" />
                        <span className="hidden font-bold sm:inline-block tracking-tight text-lg">CoreForge</span>
                    </Link>
                    <nav className="flex items-center gap-6 text-sm font-medium">
                        <Link href="/docs/components/error-handler" className="transition-colors hover:text-foreground/80 text-foreground/60">{t("docs")}</Link>
                        <Link href="/components" className="transition-colors hover:text-foreground/80 text-foreground/60">{t("components")}</Link>
                    </nav>
                </div>

                <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
                    <div className="w-full flex-1 md:w-auto md:flex-none">
                        <div className="relative">
                            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder={t("searchPlaceholder")}
                                className="h-9 md:w-[300px] lg:w-[300px] pl-9 bg-muted/40 rounded-md border-muted-foreground/20"
                            />
                            <kbd className="pointer-events-none absolute right-1.5 top-2 hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
                                <span className="text-xs">⌘</span>K
                            </kbd>
                        </div>
                    </div>
                    <nav className="flex items-center gap-1">
                        <LanguageSwitcher />
                        <Button variant="ghost" size="icon" className="w-9 px-0">
                            <Github className="h-4 w-4" />
                        </Button>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="w-9 px-0"
                            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
                        >
                            <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                            <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                            <span className="sr-only">Toggle theme</span>
                        </Button>
                    </nav>
                </div>
            </div>
        </header>
    );
}
