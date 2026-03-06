import { useTranslations } from "next-intl";
import { Link } from "@/i18n/routing";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Layers, ShieldCheck, Zap, TerminalSquare } from "lucide-react";
import { InteractiveTerminal } from "@/components/interactive-terminal";
import { TechMarquee } from "@/components/tech-marquee";

export default function Home() {
  const t = useTranslations("Hero");
  const f = useTranslations("Features");
  const hw = useTranslations("HowItWorks");
  const faq = useTranslations("FAQ");
  const cta = useTranslations("CTA");

  return (
    <div className="flex flex-col items-center justify-center pt-20 pb-16 overflow-hidden">
      {/* HERO SECTION */}
      <section className="mx-auto grid max-w-[1200px] grid-cols-1 lg:grid-cols-2 gap-12 items-center px-6 mb-16 pt-8">
        {/* Left: Text Content */}
        <div className="flex flex-col items-start text-left animate-in fade-in slide-in-from-left-4 duration-500 ease-in-out">
          <Link href="/docs/components/error-handler">
            <Badge variant="outline" className="mb-6 py-1.5 px-4 bg-orange-500/10 hover:bg-orange-500/20 border-orange-500/20 text-orange-400 cursor-pointer font-medium">
              {t("badge")}
            </Badge>
          </Link>
          <h1 className="text-4xl sm:text-5xl md:text-6xl font-extrabold tracking-tight mb-6 text-foreground/90 leading-tight">
            {t("title1")} <br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-orange-400 via-rose-500 to-purple-600">
              {t("title2")}
            </span>
          </h1>
          <p className="max-w-[550px] text-lg sm:text-xl text-muted-foreground mb-10 font-medium leading-relaxed">
            {t("subtitle")}
          </p>
          <div className="flex flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
            <Link href="/docs" className="w-full sm:w-auto">
              <Button size="lg" className="w-full sm:w-auto rounded-lg font-semibold h-12 px-8 bg-foreground text-background hover:bg-foreground/90 shadow-lg hover:shadow-orange-500/20 transition-all">
                {t("getStarted")}
              </Button>
            </Link>
            <div className="flex items-center gap-3 text-sm text-muted-foreground font-mono bg-muted/40 px-5 py-3 rounded-lg border border-muted w-full sm:w-auto justify-center hover:bg-muted/60 transition-colors">
              <TerminalSquare size={16} /> <span>npm i -g forge-cli</span>
            </div>
          </div>
        </div>

        {/* Right: Interactive Terminal */}
        <div className="flex justify-center lg:justify-end animate-in fade-in slide-in-from-right-8 duration-700 ease-in-out">
          <InteractiveTerminal />
        </div>
      </section>

      {/* TECH MARQUEE */}
      <TechMarquee />

      {/* FEATURES SECTION (Why CoreForge Grid 2x2) */}
      <section className="mx-auto w-full max-w-[1200px] px-6 mt-24 mb-20">
        <div className="text-center mb-16 animate-in fade-in slide-in-from-bottom-4 duration-500 delay-200">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">{f("title")}</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 lg:gap-8">
          <Card className="flex flex-col items-start text-left p-8 bg-white/5 border-white/10 hover:bg-white/10 hover:border-white/20 transition-all duration-300 group">
            <div className="p-3.5 bg-orange-500/10 rounded-xl text-orange-400 mb-5 group-hover:scale-110 group-hover:bg-orange-500/20 transition-all">
              <Layers className="h-7 w-7" />
            </div>
            <h3 className="font-bold text-xl mb-3 tracking-tight text-foreground/90">{f("f1_title")}</h3>
            <p className="text-muted-foreground text-base leading-relaxed">{f("f1_desc")}</p>
          </Card>

          <Card className="flex flex-col items-start text-left p-8 bg-white/5 border-white/10 hover:bg-white/10 hover:border-white/20 transition-all duration-300 group">
            <div className="p-3.5 bg-blue-500/10 rounded-xl text-blue-400 mb-5 group-hover:scale-110 group-hover:bg-blue-500/20 transition-all">
              <Zap className="h-7 w-7" />
            </div>
            <h3 className="font-bold text-xl mb-3 tracking-tight text-foreground/90">{f("f2_title")}</h3>
            <p className="text-muted-foreground text-base leading-relaxed">{f("f2_desc")}</p>
          </Card>

          <Card className="flex flex-col items-start text-left p-8 bg-white/5 border-white/10 hover:bg-white/10 hover:border-white/20 transition-all duration-300 group">
            <div className="p-3.5 bg-purple-500/10 rounded-xl text-purple-400 mb-5 group-hover:scale-110 group-hover:bg-purple-500/20 transition-all">
              <TerminalSquare className="h-7 w-7" />
            </div>
            <h3 className="font-bold text-xl mb-3 tracking-tight text-foreground/90">{f("f3_title")}</h3>
            <p className="text-muted-foreground text-base leading-relaxed">{f("f3_desc")}</p>
          </Card>

          <Card className="flex flex-col items-start text-left p-8 bg-white/5 border-white/10 hover:bg-white/10 hover:border-white/20 transition-all duration-300 group">
            <div className="p-3.5 bg-emerald-500/10 rounded-xl text-emerald-400 mb-5 group-hover:scale-110 group-hover:bg-emerald-500/20 transition-all">
              <ShieldCheck className="h-7 w-7" />
            </div>
            <h3 className="font-bold text-xl mb-3 tracking-tight text-foreground/90">{f("f4_title")}</h3>
            <p className="text-muted-foreground text-base leading-relaxed">{f("f4_desc")}</p>
          </Card>
        </div>
      </section>

      {/* HOW IT WORKS SECTION */}
      <section className="mx-auto w-full max-w-[1200px] px-6 mt-16 mb-20 bg-muted/10 py-16 rounded-[2rem] border border-white/5">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">{hw("title")}</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-12 lg:gap-8 px-4">
          {/* Step 1 */}
          <div className="flex flex-col items-center text-center relative">
            <div className="hidden md:block absolute top-8 left-[60%] w-full h-[2px] bg-gradient-to-r from-orange-500/20 to-transparent z-0"></div>
            <div className="flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-orange-500/20 to-rose-500/10 text-orange-400 font-black text-2xl mb-8 relative z-10 border border-orange-500/20 shadow-lg shadow-orange-500/5">
              1
            </div>
            <h3 className="font-bold text-xl mb-3 tracking-tight">{hw("step1_title")}</h3>
            <p className="text-muted-foreground text-base max-w-[280px]">{hw("step1_desc")}</p>
          </div>
          {/* Step 2 */}
          <div className="flex flex-col items-center text-center relative">
            <div className="hidden md:block absolute top-8 left-[60%] w-full h-[2px] bg-gradient-to-r from-orange-500/20 to-transparent z-0"></div>
            <div className="flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-orange-500/20 to-rose-500/10 text-orange-400 font-black text-2xl mb-8 relative z-10 border border-orange-500/20 shadow-lg shadow-orange-500/5">
              2
            </div>
            <h3 className="font-bold text-xl mb-3 tracking-tight">{hw("step2_title")}</h3>
            <p className="text-muted-foreground text-base max-w-[280px]">{hw("step2_desc")}</p>
          </div>
          {/* Step 3 */}
          <div className="flex flex-col items-center text-center">
            <div className="flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-orange-500/20 to-rose-500/10 text-orange-400 font-black text-2xl mb-8 relative z-10 border border-orange-500/20 shadow-lg shadow-orange-500/5">
              3
            </div>
            <h3 className="font-bold text-xl mb-3 tracking-tight">{hw("step3_title")}</h3>
            <p className="text-muted-foreground text-base max-w-[280px]">{hw("step3_desc")}</p>
          </div>
        </div>
      </section>

      {/* CTA SECTION */}
      <section className="mx-auto w-full px-4 mt-16 pb-12 text-center">
        <div className="max-w-[900px] mx-auto bg-gradient-to-br from-orange-500/10 via-rose-500/5 to-purple-500/10 rounded-[2.5rem] p-12 md:p-20 border border-orange-500/10 relative overflow-hidden group">
          {/* Subtle background glow */}
          <div className="absolute inset-x-0 bottom-0 top-1/2 bg-gradient-to-t from-orange-500/10 w-full blur-3xl opacity-50 group-hover:opacity-100 transition-opacity duration-700" />

          <div className="relative z-10">
            <h2 className="text-4xl md:text-5xl font-extrabold tracking-tight mb-6">{cta("title")}</h2>
            <p className="text-lg md:text-xl text-muted-foreground mb-10 max-w-[600px] mx-auto font-medium">
              {cta("subtitle")}
            </p>
            <Link href="/docs/installation">
              <Button size="lg" className="rounded-xl font-bold h-14 px-10 text-lg shadow-xl shadow-orange-500/20 hover:shadow-orange-500/40 bg-foreground text-background hover:bg-foreground/90 transition-all w-full sm:w-auto hover:scale-105 active:scale-95">
                {cta("action")}
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* FAQ SECTION */}
      <section className="mx-auto w-full max-w-[800px] px-6 mt-24 mb-10">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold tracking-tight mb-4">{faq("title")}</h2>
        </div>
        <div className="flex flex-col gap-6">
          <Card className="p-6 md:p-8 text-left border-white/5 bg-white/5 hover:bg-white/10 transition-colors duration-300">
            <h3 className="font-bold text-lg mb-3 flex items-start gap-3 text-foreground/90">
              <span className="text-orange-400 mt-0.5">Q.</span> {faq("q1")}
            </h3>
            <p className="text-muted-foreground text-base leading-relaxed pl-8">{faq("a1")}</p>
          </Card>
          <Card className="p-6 md:p-8 text-left border-white/5 bg-white/5 hover:bg-white/10 transition-colors duration-300">
            <h3 className="font-bold text-lg mb-3 flex items-start gap-3 text-foreground/90">
              <span className="text-orange-400 mt-0.5">Q.</span> {faq("q2")}
            </h3>
            <p className="text-muted-foreground text-base leading-relaxed pl-8">{faq("a2")}</p>
          </Card>
          <Card className="p-6 md:p-8 text-left border-white/5 bg-white/5 hover:bg-white/10 transition-colors duration-300">
            <h3 className="font-bold text-lg mb-3 flex items-start gap-3 text-foreground/90">
              <span className="text-orange-400 mt-0.5">Q.</span> {faq("q3")}
            </h3>
            <p className="text-muted-foreground text-base leading-relaxed pl-8">{faq("a3")}</p>
          </Card>
        </div>
      </section>
    </div>
  );
}
