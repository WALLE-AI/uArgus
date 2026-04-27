"use client";

import { motion } from "framer-motion";
import { User, Video, FlaskConical, Newspaper } from "lucide-react";
import { useI18n } from "@/lib/i18n";

const CASE_DEFS = [
  { icon: User, prefix: "cases.media", color: "border-accent-blue/30" },
  { icon: Video, prefix: "cases.live", color: "border-accent-green/30" },
  { icon: FlaskConical, prefix: "cases.research", color: "border-accent-purple/30" },
  { icon: Newspaper, prefix: "cases.news", color: "border-accent-orange/30" },
];

export default function UseCases() {
  const { t } = useI18n();
  return (
    <section id="cases" className="py-24 bg-muted/30">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("cases.title")}
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            {t("cases.subtitle")}
          </p>
        </div>

        <div className="grid gap-6 sm:grid-cols-2">
          {CASE_DEFS.map((c, i) => (
            <motion.div
              key={c.prefix}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-40px" }}
              transition={{ duration: 0.5, delay: i * 0.1 }}
              className={`rounded-2xl border ${c.color} bg-white p-8`}
            >
              <div className="mb-4 flex items-center gap-3">
                <div className="rounded-xl bg-muted p-2.5">
                  <c.icon className="h-5 w-5 text-muted-foreground" />
                </div>
                <h3 className="text-lg font-semibold">{t(`${c.prefix}.title`)}</h3>
              </div>

              <div className="mb-4 grid grid-cols-2 gap-3">
                <div className="rounded-lg bg-red-50 p-3">
                  <div className="text-[10px] font-medium uppercase text-red-400 mb-1">{t("cases.before")}</div>
                  <div className="text-xs text-red-600">{t(`${c.prefix}.before`)}</div>
                </div>
                <div className="rounded-lg bg-green-50 p-3">
                  <div className="text-[10px] font-medium uppercase text-green-500 mb-1">{t("cases.after")}</div>
                  <div className="text-xs text-green-700">{t(`${c.prefix}.after`)}</div>
                </div>
              </div>

              <blockquote className="text-sm italic text-muted-foreground leading-relaxed">
                &ldquo;{t(`${c.prefix}.quote`)}&rdquo;
              </blockquote>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
