"use client";

import { motion } from "framer-motion";
import { ArrowRight, Play } from "lucide-react";
import dynamic from "next/dynamic";
import { useI18n } from "@/lib/i18n";

const Globe = dynamic(() => import("./Globe"), { ssr: false });

export default function Hero() {
  const { t } = useI18n();

  const STATS = [
    { value: "3000+", label: t("hero.stat.sources") },
    { value: "50+", label: t("hero.stat.dimensions") },
    { value: "24/7", label: t("hero.stat.analysis") },
    { value: "10x", label: t("hero.stat.cost") },
  ];

  return (
    <section className="relative min-h-screen overflow-hidden bg-hero pt-16">
      {/* Gradient overlays */}
      <div className="pointer-events-none absolute inset-0">
        <div className="absolute left-1/4 top-1/4 h-96 w-96 rounded-full bg-accent-blue/10 blur-[128px]" />
        <div className="absolute right-1/4 bottom-1/4 h-96 w-96 rounded-full bg-accent-green/10 blur-[128px]" />
      </div>

      <div className="relative mx-auto grid max-w-7xl gap-12 px-6 py-24 lg:grid-cols-2 lg:py-32 items-center">
        {/* Left: Text */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, ease: "easeOut" }}
          className="flex flex-col gap-8"
        >
          <div className="inline-flex w-fit items-center gap-2 rounded-full border border-accent-blue/30 bg-accent-blue/10 px-4 py-1.5 text-xs text-accent-blue">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-accent-blue opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-accent-blue" />
            </span>
            {t("hero.badge")}
          </div>

          <h1 className="text-4xl font-bold leading-tight tracking-tight text-white sm:text-5xl lg:text-6xl">
            {t("hero.title.1")}
            <span className="text-gradient">{t("hero.title.highlight")}</span>
            {t("hero.title.2")}
          </h1>

          <p className="max-w-lg text-lg leading-relaxed text-white/60">
            {t("hero.subtitle")}
          </p>

          <div className="flex flex-wrap gap-4">
            <a
              href="#pricing"
              className="group inline-flex items-center gap-2 rounded-lg bg-accent-blue px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-accent-blue/25 transition-all hover:bg-accent-blue/90 hover:shadow-xl hover:shadow-accent-blue/30"
            >
              {t("hero.cta.trial")}
              <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
            </a>
            <button className="inline-flex items-center gap-2 rounded-lg border border-white/20 px-6 py-3 text-sm font-medium text-white/80 transition-all hover:border-white/40 hover:text-white">
              <Play className="h-4 w-4" />
              {t("hero.cta.demo")}
            </button>
          </div>
        </motion.div>

        {/* Right: Globe */}
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 1, delay: 0.3, ease: "easeOut" }}
          className="hidden lg:block"
        >
          <Globe />
        </motion.div>
      </div>

      {/* Stats ticker */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8, delay: 0.6 }}
        className="relative border-t border-white/10 bg-white/[0.02]"
      >
        <div className="mx-auto grid max-w-5xl grid-cols-2 gap-4 px-6 py-8 sm:grid-cols-4">
          {STATS.map((stat) => (
            <div key={stat.label} className="text-center">
              <div className="text-2xl font-bold text-white sm:text-3xl">{stat.value}</div>
              <div className="mt-1 text-xs text-white/50 sm:text-sm">{stat.label}</div>
            </div>
          ))}
        </div>
      </motion.div>
    </section>
  );
}
