"use client";

import { motion } from "framer-motion";
import { Radar, Brain, BotMessageSquare, Share2, BarChart3 } from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { useI18n } from "@/lib/i18n";

interface FeatureDef {
  icon: LucideIcon;
  titleKey: string;
  descKey: string;
  color: string;
  bgColor: string;
}

const FEATURE_DEFS: FeatureDef[] = [
  { icon: Radar, titleKey: "features.monitoring.title", descKey: "features.monitoring.desc", color: "text-accent-blue", bgColor: "bg-accent-blue/10" },
  { icon: BarChart3, titleKey: "features.analysis.title", descKey: "features.analysis.desc", color: "text-accent-green", bgColor: "bg-accent-green/10" },
  { icon: Brain, titleKey: "features.knowledge.title", descKey: "features.knowledge.desc", color: "text-accent-purple", bgColor: "bg-accent-purple/10" },
  { icon: BotMessageSquare, titleKey: "features.assistant.title", descKey: "features.assistant.desc", color: "text-accent-orange", bgColor: "bg-accent-orange/10" },
  { icon: Share2, titleKey: "features.distribution.title", descKey: "features.distribution.desc", color: "text-accent-blue", bgColor: "bg-accent-blue/10" },
];

const cardVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { duration: 0.5, delay: i * 0.1, ease: "easeOut" as const },
  }),
};

export default function FeatureCards() {
  const { t } = useI18n();
  return (
    <section id="features" className="py-24 bg-white">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("features.title.1")}
            <span className="text-gradient">{t("features.title.highlight")}</span>
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("features.subtitle")}
          </p>
        </div>

        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {FEATURE_DEFS.map((feature, i) => (
            <motion.div
              key={feature.titleKey}
              custom={i}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true, margin: "-50px" }}
              variants={cardVariants}
              className="group relative rounded-2xl border border-border bg-white p-8 transition-all hover:border-accent-blue/30 hover:shadow-lg hover:shadow-accent-blue/5"
            >
              <div className={`mb-4 inline-flex rounded-xl ${feature.bgColor} p-3`}>
                <feature.icon className={`h-6 w-6 ${feature.color}`} />
              </div>
              <h3 className="mb-2 text-lg font-semibold">{t(feature.titleKey)}</h3>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {t(feature.descKey)}
              </p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
