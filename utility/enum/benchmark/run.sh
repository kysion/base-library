#!/bin/bash

# 执行所有基准测试并输出结果
cd "$(dirname "$0")"
go test -bench=. -benchmem > benchmark_results.txt

# 分析结果并输出摘要
echo -e "\n性能测试结果摘要:" | tee benchmark_results.txt
echo "====================================" | tee -a benchmark_results.txt

# 显示测试结果
cat benchmark_results.txt

echo -e "\n详细结果保存在 benchmark_results.txt 文件中" 