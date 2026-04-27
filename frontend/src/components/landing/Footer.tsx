"use client";

import { Eye } from "lucide-react";
import { useI18n } from "@/lib/i18n";

export default function Footer() {
  const { t, toggle } = useI18n();

  const LINKS: Record<string, string[]> = {
    [t("footer.product")]: ["Features", "Pricing", "Changelog", "Docs"],
    [t("footer.company")]: ["About", "Blog", "Careers", "Contact"],
    [t("footer.legal")]: ["Privacy", "Terms", "Security"],
    [t("footer.connect")]: ["Twitter", "GitHub", "Discord"],
  };

  return (
    <footer className="border-t border-border bg-muted/30">
      <div className="mx-auto max-w-7xl px-6 py-16">
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-5">
          {/* Brand */}
          <div className="lg:col-span-2">
            <a href="#" className="flex items-center gap-2 font-bold text-lg">
              <Eye className="h-5 w-5 text-accent-blue" />
              uArgus
            </a>
            <p className="mt-3 max-w-xs text-sm text-muted-foreground leading-relaxed">
              {t("footer.desc")}
            </p>
          </div>

          {/* Link columns */}
          {Object.entries(LINKS).map(([title, items]) => (
            <div key={title}>
              <h4 className="mb-3 text-sm font-semibold">{title}</h4>
              <ul className="space-y-2">
                {items.map((item) => (
                  <li key={item}>
                    <a href="#" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                      {item}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="mt-12 flex flex-col items-center justify-between gap-4 border-t border-border pt-8 sm:flex-row">
          <p className="text-xs text-muted-foreground">
            &copy; {new Date().getFullYear()} {t("footer.copyright")}
          </p>
          <div className="flex gap-4">
            <button
              onClick={toggle}
              className="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              {t("nav.lang")}
            </button>
          </div>
        </div>
      </div>
    </footer>
  );
}
