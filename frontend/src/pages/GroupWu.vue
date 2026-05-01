<script setup>
import { ref, computed, onMounted } from 'vue'
import { NCard, NButton, NSpace, NTable, NStatistic, NGrid, NGi, NEmpty, useMessage } from 'naive-ui'
import { GetGroupWu } from '../../wailsjs/go/main/App'
import { RefreshCw } from 'lucide-vue-next'

const nmessage = useMessage()
const groupdata = ref([])
const loading = ref(false)

const totalMembers = computed(() => groupdata.value.reduce((sum, g) => sum + g.member_count, 0))
const totalWu = computed(() => groupdata.value.reduce((sum, g) => sum + g.total_wu, 0))
const avgWu = computed(() => groupdata.value.length > 0 ? Math.round(totalWu.value / groupdata.value.length) : 0)

function formatWu(val) {
    const n = Math.floor(val)
    if (n >= 10000) {
        return (n / 10000).toFixed(2) + '万'
    }
    return n
}

function getData() {
    loading.value = true
    groupdata.value = []
    GetGroupWu().then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            groupdata.value = resp.data
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(() => {
        nmessage.error('获取分组武勋数据失败')
    }).finally(() => {
        loading.value = false
    })
}

onMounted(() => {
    getData()
})
</script>

<template>
    <div class="page-groupwu">
        <n-grid :cols="3" :x-gap="16" :y-gap="16" class="stat-grid">
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="总人数" :value="totalMembers" />
                </n-card>
            </n-gi>
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="总武勋" :value="formatWu(totalWu)" />
                </n-card>
            </n-gi>
            <n-gi>
                <n-card embedded size="small">
                    <n-statistic label="平均武勋" :value="formatWu(avgWu)" />
                </n-card>
            </n-gi>
        </n-grid>

        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">分组武勋</h2>
                    <p class="page-desc">更新武勋数据请同步成员数据</p>
                </div>
                <n-space>
                    <n-button @click="getData" :loading="loading">
                        <template #icon><RefreshCw :size="16" /></template>
                        刷新
                    </n-button>
                </n-space>
            </div>

            <n-table v-if="groupdata.length > 0" :bordered="true" :single-line="false" class="styled-table">
                <thead>
                    <tr>
                        <th>分组名称</th>
                        <th>人数</th>
                        <th>总武勋</th>
                        <th>平均武勋</th>
                        <th>0武勋人数</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="u in groupdata" :key="u.group">
                        <td>
                            <span class="group-name">{{ u.group }}</span>
                        </td>
                        <td>{{ u.member_count }}</td>
                        <td>{{ formatWu(u.total_wu) }}</td>
                        <td>{{ formatWu(u.average_wu) }}</td>
                        <td>
                            <span :class="{ 'zero-warn': u.zero_wu_count > 0 }">{{ u.zero_wu_count }}</span>
                        </td>
                    </tr>
                </tbody>
            </n-table>
            <n-empty v-else description="暂无分组武勋数据" style="padding: 60px 0;" />
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-groupwu {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.page-card {
    border-radius: 12px;
}

.page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 20px;
}

.page-title {
    font-size: 20px;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 4px;
}

.page-desc {
    font-size: 13px;
    color: #64748b;
}

.styled-table {
    border-radius: 8px;
    overflow: hidden;

    thead {
        th {
            background: #f8f9fb;
            font-weight: 600;
            color: #64748b;
            font-size: 13px;
            padding: 12px 16px;
        }
    }

    tbody {
        td {
            padding: 12px 16px;
            font-size: 14px;
        }

        tr:nth-child(even) {
            background: #fafbfc;
        }

        tr:hover td {
            background: #f1f3f5;
        }
    }
}

.group-name {
    font-weight: 600;
    color: #1e293b;
}

.zero-warn {
    color: #ef4444;
    font-weight: 600;
}
</style>
