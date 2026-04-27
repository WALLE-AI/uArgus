"use client";

import { motion } from "framer-motion";
import { useI18n } from "@/lib/i18n";

const SOURCES = [
  "RSS Feeds", "arXiv", "HackerNews", "Reddit", "GitHub",
  "Reuters", "Bloomberg", "TechCrunch", "Nature", "PubMed",
  "USGS", "NASA", "IMF", "World Bank", "CoinGecko",
];

export default function TrustBar() {
  const { t } = useI18n();
  return (
    <section className="border-y border-border bg-muted/50 py-8 overflow-hidden">
      <div className="mx-auto max-w-7xl px-6">
        <p className="mb-6 text-center text-xs font-medium uppercase tracking-widest text-muted-foreground">
          {t("trust.title")}
        </p>
        <motion.div
          className="flex gap-8"
          animate={{ x: [0, -1200] }}
          transition={{ duration: 30, repeat: Infinity, ease: "linear" }}
        >
          {[...SOURCES, ...SOURCES].map((name, i) => (
            <div
              key={`${name}-${i}`}
              className="flex shrink-0 items-center gap-2 rounded-full border border-border bg-white px-4 py-2 text-sm text-muted-foreground"
            >
              <div className="h-2 w-2 rounded-full bg-accent-green/60" />
              {name}
            </div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}
