import { I18nProvider } from "@/lib/i18n";
import Navbar from "@/components/landing/Navbar";
import Hero from "@/components/landing/Hero";
import TrustBar from "@/components/landing/TrustBar";
import FeatureCards from "@/components/landing/FeatureCards";
import HowItWorks from "@/components/landing/HowItWorks";
import Screenshots from "@/components/landing/Screenshots";
import ROINumbers from "@/components/landing/ROINumbers";
import UseCases from "@/components/landing/UseCases";
import Pricing from "@/components/landing/Pricing";
import Footer from "@/components/landing/Footer";

export default function Home() {
  return (
    <I18nProvider>
      <Navbar />
      <Hero />
      <TrustBar />
      <FeatureCards />
      <HowItWorks />
      <Screenshots />
      <ROINumbers />
      <UseCases />
      <Pricing />
      <Footer />
    </I18nProvider>
  );
}
