"use client";

import { useState } from "react";
import { Menu, X, Eye } from "lucide-react";
import { cn } from "@/lib/utils";
import { useI18n } from "@/lib/i18n";

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const { t, toggle } = useI18n();

  const NAV_LINKS = [
    { label: t("nav.features"), href: "#features" },
    { label: t("nav.howItWorks"), href: "#how-it-works" },
    { label: t("nav.cases"), href: "#cases" },
    { label: t("nav.pricing"), href: "#pricing" },
  ];

  return (
    <nav className="fixed top-0 z-50 w-full border-b border-white/10 bg-hero-from/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
        {/* Logo */}
        <a href="#" className="flex items-center gap-2 text-white font-bold text-xl">
          <Eye className="h-6 w-6 text-accent-blue" />
          <span>uArgus</span>
        </a>

        {/* Desktop nav */}
        <div className="hidden items-center gap-8 md:flex">
          {NAV_LINKS.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="text-sm text-white/70 transition-colors hover:text-white"
            >
              {link.label}
            </a>
          ))}
        </div>

        {/* CTA + lang */}
        <div className="hidden items-center gap-4 md:flex">
          <button
            onClick={toggle}
            className="text-sm text-white/70 hover:text-white transition-colors"
          >
            {t("nav.lang")}
          </button>
          <a
            href="#pricing"
            className="rounded-lg bg-accent-blue px-4 py-2 text-sm font-medium text-white transition-all hover:bg-accent-blue/90 hover:shadow-lg hover:shadow-accent-blue/25"
          >
            {t("nav.startFree")}
          </a>
        </div>

        {/* Mobile hamburger */}
        <button
          className="text-white md:hidden"
          onClick={() => setOpen(!open)}
          aria-label="Toggle menu"
        >
          {open ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
        </button>
      </div>

      {/* Mobile menu */}
      <div
        className={cn(
          "overflow-hidden transition-all duration-300 md:hidden",
          open ? "max-h-80" : "max-h-0"
        )}
      >
        <div className="flex flex-col gap-4 border-t border-white/10 bg-hero-from/95 px-6 py-6">
          {NAV_LINKS.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="text-sm text-white/70 hover:text-white"
              onClick={() => setOpen(false)}
            >
              {link.label}
            </a>
          ))}
          <button
            onClick={toggle}
            className="text-sm text-white/70 hover:text-white text-left"
          >
            {t("nav.lang")}
          </button>
          <a
            href="#pricing"
            className="mt-2 rounded-lg bg-accent-blue px-4 py-2 text-center text-sm font-medium text-white"
          >
            {t("nav.startFree")}
          </a>
        </div>
      </div>
    </nav>
  );
}
