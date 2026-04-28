"use client";

import AreaChartCard from "@/components/charts/AreaChartCard";
import BarChartCard from "@/components/charts/BarChartCard";
import LineChartCard from "@/components/charts/LineChartCard";

const CHANNELS = [
  { platform: "微信公众号", views: "45.2K", followers: "+342", engagement: "4.2%", revenue: "¥3,240", color: "#07C160" },
  { platform: "知乎", views: "28.1K", followers: "+156", engagement: "6.8%", revenue: "¥1,820", color: "#0066FF" },
  { platform: "小红书", views: "12.5K", followers: "+89", engagement: "8.3%", revenue: "¥960", color: "#FE2C55" },
  { platform: "抖音", views: "67.8K", followers: "+1.2K", engagement: "3.1%", revenue: "¥5,100", color: "#000000" },
];

const TRAFFIC_DATA = Array.from({ length: 30 }, (_, i) => ({
  date: `${i + 1}`,
  微信: Math.round(1200 + Math.sin(i / 3) * 400 + Math.random() * 200),
  知乎: Math.round(800 + Math.cos(i / 4) * 300 + Math.random() * 150),
  小红书: Math.round(400 + Math.sin(i / 5) * 200 + Math.random() * 100),
  抖音: Math.round(2000 + Math.sin(i / 2) * 600 + Math.random() * 400),
}));

const FOLLOWER_DATA = Array.from({ length: 12 }, (_, i) => ({
  month: `${i + 1}月`,
  微信: Math.round(8000 + i * 350 + Math.random() * 100),
  知乎: Math.round(5000 + i * 200 + Math.random() * 80),
  小红书: Math.round(2000 + i * 150 + Math.random() * 50),
  抖音: Math.round(3000 + i * 500 + Math.random() * 200),
}));

const ENGAGEMENT_DATA = [
  { platform: "微信", likes: 1240, comments: 342, shares: 189 },
  { platform: "知乎", likes: 2100, comments: 876, shares: 145 },
  { platform: "小红书", likes: 890, comments: 234, shares: 432 },
  { platform: "抖音", likes: 5600, comments: 1200, shares: 890 },
];

export default function ChannelsPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">渠道运营</h1>

      {/* Platform metric cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {CHANNELS.map((ch) => (
          <div key={ch.platform} className="rounded-xl border border-border bg-white p-5">
            <div className="flex items-center justify-between">
              <div className="text-sm font-semibold">{ch.platform}</div>
              <div className="h-3 w-3 rounded-full" style={{ backgroundColor: ch.color }} />
            </div>
            <div className="mt-3 space-y-1.5">
              {[
                { label: "阅读量", value: ch.views },
                { label: "新增粉丝", value: ch.followers, cls: "text-accent-green" },
                { label: "互动率", value: ch.engagement },
                { label: "收入", value: ch.revenue, cls: "text-accent-blue" },
              ].map((row) => (
                <div key={row.label} className="flex justify-between text-xs">
                  <span className="text-muted-foreground">{row.label}</span>
                  <span className={`font-medium ${row.cls ?? ""}`}>{row.value}</span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* Traffic trend */}
      <LineChartCard
        title="流量趋势（30天）"
        subtitle="各平台每日阅读量"
        data={TRAFFIC_DATA}
        xKey="date"
        lines={[
          { key: "微信", color: "#07C160" },
          { key: "知乎", color: "#0066FF" },
          { key: "小红书", color: "#FE2C55" },
          { key: "抖音", color: "#000000" },
        ]}
        height={280}
      />

      {/* Follower growth + engagement */}
      <div className="grid gap-4 lg:grid-cols-2">
        <AreaChartCard
          title="粉丝增长趋势（12个月）"
          data={FOLLOWER_DATA}
          xKey="month"
          areas={[
            { key: "微信", color: "#07C160" },
            { key: "知乎", color: "#0066FF" },
            { key: "小红书", color: "#FE2C55" },
            { key: "抖音", color: "#000000" },
          ]}
          stacked
          height={260}
        />
        <BarChartCard
          title="互动指标对比"
          subtitle="点赞 / 评论 / 分享"
          data={ENGAGEMENT_DATA}
          xKey="platform"
          bars={[
            { key: "likes", color: "#EF4444", name: "点赞" },
            { key: "comments", color: "#3B82F6", name: "评论" },
            { key: "shares", color: "#10B981", name: "分享" },
          ]}
          height={260}
        />
      </div>

      {/* ROI table */}
      <div className="rounded-xl border border-border bg-white p-6">
        <h3 className="text-sm font-semibold mb-4">ROI 汇总</h3>
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-left text-xs text-muted-foreground">
              <th className="pb-2 font-medium">平台</th>
              <th className="pb-2 font-medium">投入成本</th>
              <th className="pb-2 font-medium">收入</th>
              <th className="pb-2 font-medium">ROI</th>
              <th className="pb-2 font-medium text-right">单粉成本</th>
            </tr>
          </thead>
          <tbody>
            {[
              { platform: "微信公众号", cost: "¥1,200", revenue: "¥3,240", roi: "170%", cpa: "¥3.51" },
              { platform: "知乎", cost: "¥800", revenue: "¥1,820", roi: "127%", cpa: "¥5.13" },
              { platform: "小红书", cost: "¥500", revenue: "¥960", roi: "92%", cpa: "¥5.62" },
              { platform: "抖音", cost: "¥2,000", revenue: "¥5,100", roi: "155%", cpa: "¥1.67" },
            ].map((r) => (
              <tr key={r.platform} className="border-b border-border last:border-0">
                <td className="py-2 font-medium">{r.platform}</td>
                <td className="py-2 text-muted-foreground">{r.cost}</td>
                <td className="py-2 text-accent-green font-medium">{r.revenue}</td>
                <td className="py-2 font-medium">{r.roi}</td>
                <td className="py-2 text-right text-muted-foreground">{r.cpa}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
