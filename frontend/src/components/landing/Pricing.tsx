"use client";

import { motion } from "framer-motion";
import { Check } from "lucide-react";
import { cn } from "@/lib/utils";
import { useI18n } from "@/lib/i18n";

interface PlanDef {
  prefix: string;
  featureCount: number;
  popular?: boolean;
}

const PLAN_DEFS: PlanDef[] = [
  { prefix: "pricing.free", featureCount: 4 },
  { prefix: "pricing.pro", featureCount: 6, popular: true },
  { prefix: "pricing.enterprise", featureCount: 7 },
];

export default function Pricing() {
  const { t } = useI18n();
  return (
    <section id="pricing" className="py-24 bg-white">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("pricing.title")}
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("pricing.subtitle")}
          </p>
        </div>

        <div className="grid gap-6 sm:grid-cols-3 max-w-5xl mx-auto">
          {PLAN_DEFS.map((plan, i) => {
            const period = t(`${plan.prefix}.period`);
            const features = Array.from({ length: plan.featureCount }, (_, j) =>
              t(`${plan.prefix}.f${j + 1}`)
            );
            return (
              <motion.div
                key={plan.prefix}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.5, delay: i * 0.1 }}
                className={cn(
                  "relative rounded-2xl border p-8 flex flex-col",
                  plan.popular
                    ? "border-accent-blue bg-hero-from text-white shadow-xl shadow-accent-blue/10 scale-105"
                    : "border-border bg-white"
                )}
              >
                {plan.popular && (
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-accent-blue px-3 py-1 text-xs font-medium text-white">
                    {t("pricing.popular")}
                  </div>
                )}

                <h3 className={cn("text-lg font-semibold", plan.popular ? "text-white" : "")}>{t(plan.prefix)}</h3>
                <div className="mt-4 flex items-baseline gap-1">
                  <span className={cn("text-4xl font-bold", plan.popular ? "text-white" : "")}>{t(`${plan.prefix}.price`)}</span>
                  {period && (
                    <span className={cn("text-sm", plan.popular ? "text-white/60" : "text-muted-foreground")}>
                      {period}
                    </span>
                  )}
                </div>
                <p className={cn("mt-2 text-sm", plan.popular ? "text-white/60" : "text-muted-foreground")}>
                  {t(`${plan.prefix}.desc`)}
                </p>

                <ul className="mt-6 flex-1 space-y-3">
                  {features.map((f) => (
                    <li key={f} className="flex items-start gap-2 text-sm">
                      <Check className={cn("h-4 w-4 shrink-0 mt-0.5 text-accent-green")} />
                      <span className={plan.popular ? "text-white/80" : "text-muted-foreground"}>{f}</span>
                    </li>
                  ))}
                </ul>

                <button
                  className={cn(
                    "mt-8 w-full rounded-lg py-2.5 text-sm font-medium transition-all",
                    plan.popular
                      ? "bg-accent-blue text-white hover:bg-accent-blue/90 shadow-lg shadow-accent-blue/25"
                      : "border border-border bg-muted hover:bg-muted/80"
                  )}
                >
                  {t(`${plan.prefix}.cta`)}
                </button>
              </motion.div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
