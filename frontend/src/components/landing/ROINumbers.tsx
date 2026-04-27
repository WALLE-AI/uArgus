"use client";

import { useEffect, useRef, useState } from "react";
import { motion, useInView } from "framer-motion";
import { TrendingDown, TrendingUp, Database as DbIcon } from "lucide-react";
import { useI18n } from "@/lib/i18n";

interface MetricDef {
  icon: React.ElementType;
  value: number;
  suffix: string;
  labelKey: string;
  descKey: string;
}

const METRIC_DEFS: MetricDef[] = [
  { icon: TrendingDown, value: 10, suffix: "x", labelKey: "roi.cost.label", descKey: "roi.cost.desc" },
  { icon: TrendingUp, value: 10, suffix: "x", labelKey: "roi.exposure.label", descKey: "roi.exposure.desc" },
  { icon: DbIcon, value: 100, suffix: "%", labelKey: "roi.retention.label", descKey: "roi.retention.desc" },
];

function AnimatedCounter({ target, suffix }: { target: number; suffix: string }) {
  const [count, setCount] = useState(0);
  const ref = useRef<HTMLSpanElement>(null);
  const isInView = useInView(ref, { once: true });

  useEffect(() => {
    if (!isInView) return;
    let start = 0;
    const duration = 1500;
    const step = (ts: number) => {
      if (!start) start = ts;
      const progress = Math.min((ts - start) / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3);
      setCount(Math.round(eased * target));
      if (progress < 1) requestAnimationFrame(step);
    };
    requestAnimationFrame(step);
  }, [isInView, target]);

  return (
    <span ref={ref} className="tabular-nums">
      {count}{suffix}
    </span>
  );
}

export default function ROINumbers() {
  const { t } = useI18n();
  return (
    <section className="py-24 bg-hero">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
            {t("roi.title")}
          </h2>
          <p className="mt-4 text-lg text-white/60">
            {t("roi.subtitle")}
          </p>
        </div>

        <div className="grid gap-8 sm:grid-cols-3">
          {METRIC_DEFS.map((m, i) => (
            <motion.div
              key={m.labelKey}
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6, delay: i * 0.15 }}
              className="text-center"
            >
              <div className="mx-auto mb-4 inline-flex rounded-xl bg-white/5 p-3">
                <m.icon className="h-6 w-6 text-accent-green" />
              </div>
              <div className="text-5xl font-bold text-white sm:text-6xl">
                <AnimatedCounter target={m.value} suffix={m.suffix} />
              </div>
              <div className="mt-2 text-lg font-semibold text-white">{t(m.labelKey)}</div>
              <p className="mt-2 text-sm text-white/50">{t(m.descKey)}</p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
