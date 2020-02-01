import React, { useEffect, useState } from "react";
import "./App.css";

import { LineChart, XAxis, YAxis, Line, Legend } from "recharts";

import Layout from "./components/Layout";
import { useToken } from "./state/token";

const getStars = (repos: string[], token: string) =>
  fetch(`http://localhost:3000/api/data?repos=${repos.join(",")}`, {
    headers: {
      Authorization: token
    }
  });

interface StarHistoryEntry {
  t: number;
  v: number;
}

const App = () => {
  const token = useToken();

  const [newRepo, setNewRepo] = useState("");
  const [data, setData] = useState<{ [k: string]: StarHistoryEntry[] }>({});
  const mergeData = (o: any) => setData({...data, ...o});
  const fetchStars = () => {
    if (token === null) {
      return;
    }
    getStars([newRepo], token)
      .then(r => r.json())
      .then(mergeData);
    setNewRepo("");
  };

  useEffect(() => console.log(data), [data]);

  const all = Object.keys(data)
    .map(k => data[k])
    .reduce((a, b) => [...a, ...b], []);
  const times: number[] = all.map(d => d.t);
  let [min, max] = [Math.min(...times), Math.max(...times)];
  let startYear = new Date(1000 * min).getFullYear();
  let endYear = new Date(1000 * max).getFullYear();
  let years = Array.from(
    { length: 1 + endYear - startYear },
    (_, i) => i + startYear
  );

  return (
    <Layout>
      <form onSubmit={e => {e.preventDefault();fetchStars()}}>
      <input type="text" value={newRepo} onChange={e => setNewRepo(e.target.value)} placeholder="repository" />
      <button disabled={!token}>
        Fetch stars
      </button>
      </form>
      {Object.keys(data).length > 0 && (
        <LineChart width={800} height={400} data={all}>
          <XAxis
            dataKey="t"
            domain={["dataMin", "dataMax"]}
            name="Time"
            scale="time"
            type="number"
            interval={"preserveStartEnd"}
            ticks={years.map(y => new Date(y, 0).getTime() / 1000)}
            tickFormatter={unixTime => new Date(1000 * unixTime).getFullYear()}
          />
          <YAxis
            tickFormatter={s => (s === 0 ? "0" : Math.floor(s / 1000) + "k")}
          />
          <Legend />
          {Object.keys(data).map(k => (
            <Line
              key={k}
              dataKey={"v"}
              name={k}
              data={data[k].sort((a, b) => a.t - b.t)}
            />
          ))}
        </LineChart>
      )}
    </Layout>
  );
};

export default App;
