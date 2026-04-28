"use client";

import AreaChartCard from "@/components/charts/AreaChartCard";
import BarChartCard from "@/components/charts/BarChartCard";

const PREDICTION_DATA = Array.from({ length: 30 }, (_, i) => {
  const actual = 50 + Math.sin(i / 4) * 20 + (i > 20 ? (i - 20) * 2 : 0);
  const predicted = 50 + Math.sin(i / 4) * 18 + (i > 20 ? (i - 20) * 2.2 : 0);
  return {
    date: `${i + 1}日`,
    实际值: i <= 25 ? Math.round(actual) : undefined,
    预测值: Math.round(predicted),
    上界: Math.round(predicted + 8 + Math.random() * 4),
    下界: Math.round(predicted - 8 - Math.random() * 4),
  };
});

const ACCURACY_DATA = [
  { month: "1月", accuracy: 72 },
  { month: "2月", accuracy: 75 },
  { month: "3月", accuracy: 71 },
  { month: "4月", accuracy: 78 },
  { month: "5月", accuracy: 82 },
  { month: "6月", accuracy: 80 },
  { month: "7月", accuracy: 76 },
  { month: "8月", accuracy: 83 },
  { month: "9月", accuracy: 79 },
  { month: "10月", accuracy: 81 },
  { month: "11月", accuracy: 78 },
  { month: "12月", accuracy: 84 },
];

const MODELS = [
  { name: "话题热度预测", domain: "科技", accuracy: "84.2%", lastRun: "10 min ago", status: "运行中" },
  { name: "情感趋势预测", domain: "全局", accuracy: "79.5%", lastRun: "15 min ago", status: "运行中" },
  { name: "地缘风险评估", domain: "地缘政治", accuracy: "76.8%", lastRun: "30 min ago", status: "运行中" },
  { name: "金融市场情绪", domain: "金融", accuracy: "81.3%", lastRun: "1 hr ago", status: "运行中" },
  { name: "能源价格走势", domain: "能源", accuracy: "73.1%", lastRun: "2 hr ago", status: "等待" },
  { name: "舆情爆发预警", domain: "全局", accuracy: "82.7%", lastRun: "5 min ago", status: "运行中" },
];

export default function PredictionPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">预测分析</h1>

      {/* Summary cards */}
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="rounded-xl border border-border bg-white p-5">
          <div className="text-3xl font-bold text-accent-green">78.3%</div>
          <div className="text-sm text-muted-foreground">90天平均准确率</div>
        </div>
        <div className="rounded-xl border border-border bg-white p-5">
          <div className="text-3xl font-bold">6</div>
          <div className="text-sm text-muted-foreground">活跃预测模型</div>
        </div>
        <div className="rounded-xl border border-border bg-white p-5">
          <div className="text-3xl font-bold text-accent-blue">142</div>
          <div className="text-sm text-muted-foreground">今日预测数量</div>
        </div>
      </div>

      {/* Prediction with confidence band */}
      <AreaChartCard
        title="趋势预测（置信区间）"
        subtitle="实际值 vs 预测值，含 95% 置信区间"
        data={PREDICTION_DATA}
        xKey="date"
        areas={[
          { key: "上界", color: "#3B82F6", name: "置信上界" },
          { key: "预测值", color: "#3B82F6", name: "预测值" },
          { key: "下界", color: "#3B82F6", name: "置信下界" },
        ]}
        height={300}
      />

      {/* Accuracy + models */}
      <div className="grid gap-4 lg:grid-cols-2">
        <BarChartCard
          title="月度预测准确率"
          subtitle="过去 12 个月的预测准确率追踪"
          data={ACCURACY_DATA}
          xKey="month"
          bars={[{ key: "accuracy", color: "#10B981", name: "准确率 %" }]}
          height={260}
        />

        <div className="rounded-xl border border-border bg-white p-6">
          <h3 className="text-sm font-semibold mb-4">活跃预测模型</h3>
          <div className="space-y-3">
            {MODELS.map((m) => (
              <div key={m.name} className="flex items-center gap-3 rounded-lg p-2 hover:bg-muted/30 transition-colors">
                <div className={`h-2 w-2 rounded-full shrink-0 ${m.status === "运行中" ? "bg-accent-green animate-pulse" : "bg-muted-foreground/40"}`} />
                <div className="flex-1 min-w-0">
                  <div className="text-xs font-medium truncate">{m.name}</div>
                  <div className="text-[10px] text-muted-foreground">{m.domain} · {m.lastRun}</div>
                </div>
                <div className="shrink-0 text-right">
                  <div className="text-xs font-bold text-accent-green">{m.accuracy}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
