"use client";

import { motion } from "framer-motion";
import { Database, Cpu, PenTool, Megaphone, ArrowRight } from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { useI18n } from "@/lib/i18n";

interface StepDef {
  icon: LucideIcon;
  titleKey: string;
  descKey: string;
  color: string;
}

const STEP_DEFS: StepDef[] = [
  { icon: Database, titleKey: "how.collect.title", descKey: "how.collect.desc", color: "text-accent-blue" },
  { icon: Cpu, titleKey: "how.analyze.title", descKey: "how.analyze.desc", color: "text-accent-green" },
  { icon: PenTool, titleKey: "how.generate.title", descKey: "how.generate.desc", color: "text-accent-purple" },
  { icon: Megaphone, titleKey: "how.distribute.title", descKey: "how.distribute.desc", color: "text-accent-orange" },
];

export default function HowItWorks() {
  const { t } = useI18n();
  return (
    <section id="how-it-works" className="py-24 bg-muted/30">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("how.title")}
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("how.subtitle")}
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {STEP_DEFS.map((step, i) => (
            <motion.div
              key={step.titleKey}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-40px" }}
              transition={{ duration: 0.5, delay: i * 0.15 }}
              className="relative"
            >
              <div className="flex flex-col items-center rounded-2xl border border-border bg-white p-8 text-center h-full">
                <div className="mb-1 text-5xl font-bold text-muted-foreground/20">
                  {String(i + 1).padStart(2, "0")}
                </div>
                <div className={`mb-4 inline-flex rounded-xl bg-muted p-3`}>
                  <step.icon className={`h-6 w-6 ${step.color}`} />
                </div>
                <h3 className="mb-2 text-lg font-semibold">{t(step.titleKey)}</h3>
                <p className="text-sm text-muted-foreground">{t(step.descKey)}</p>
              </div>
              {i < STEP_DEFS.length - 1 && (
                <div className="hidden lg:flex absolute -right-4 top-1/2 z-10 -translate-y-1/2 text-muted-foreground/30">
                  <ArrowRight className="h-6 w-6" />
                </div>
              )}
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
