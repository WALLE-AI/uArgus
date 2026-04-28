export default function HealthPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">源健康状态</h1>
      <div className="grid gap-4 sm:grid-cols-3">
        {[
          { label: "正常", count: 3102, pct: "95.5%", color: "border-accent-green", text: "text-accent-green" },
          { label: "降级", count: 108, pct: "3.3%", color: "border-accent-orange", text: "text-accent-orange" },
          { label: "失败", count: 37, pct: "1.1%", color: "border-red-500", text: "text-red-500" },
        ].map((s) => (
          <div key={s.label} className={`rounded-xl border-l-4 ${s.color} border border-border bg-white p-6`}>
            <div className="text-3xl font-bold">{s.count}</div>
            <div className={`text-sm font-medium ${s.text}`}>{s.label}</div>
            <div className="text-xs text-muted-foreground">{s.pct} of total</div>
          </div>
        ))}
      </div>
      <div className="rounded-xl border border-border bg-white p-6">
        <h3 className="text-sm font-semibold mb-4">健康趋势（过去 7 天）</h3>
        <div className="h-48 flex items-end gap-1">
          {[97, 96, 95, 96, 95, 94, 96].map((v, i) => (
            <div key={i} className="flex-1 flex flex-col items-center gap-1">
              <div className="w-full rounded-t bg-accent-green/60" style={{ height: `${v}%` }} />
              <span className="text-[10px] text-muted-foreground">{["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"][i]}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
