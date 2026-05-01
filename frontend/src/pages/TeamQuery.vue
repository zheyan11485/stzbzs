<script setup>
import { ref } from 'vue'
import { NCard, NButton, NInput, NEmpty, NSpin, NTag, useMessage } from 'naive-ui'
import { GetPlayerTeam } from '../../wailsjs/go/main/App'
import { Search, Swords, Star } from 'lucide-vue-next'
import { herocfg, skillcfg, gear_feature_cfg, gear_cfg } from '../cfg'

const heroMap = JSON.parse(herocfg)
const skillMap = JSON.parse(skillcfg)
const gearFeatureMap = gear_feature_cfg
const gearMap = {}
gear_cfg.forEach(g => { gearMap[String(g.gear_id)] = g })

const nmessage = useMessage()
const loading = ref(false)
const results = ref([])

const searchName = ref('')
const searchUnion = ref('')
const searchIdu = ref('')

const hasSearched = ref(false)

const doSearch = () => {
    loading.value = true
    results.value = []
    hasSearched.value = true
    GetPlayerTeam(searchName.value, searchUnion.value, searchIdu.value).then(v => {
        let resp = JSON.parse(v)
        if (resp.code == 200) {
            results.value = resp.data || []
        } else {
            nmessage.error(resp.msg)
        }
    }).catch(e => {
        nmessage.error('查询失败: ' + e)
    }).finally(() => {
        loading.value = false
    })
}

const resolveHeroId = (id) => {
    if (!id) return id
    const num = Number(id)
    return num >= 130000 ? num - 30000 : num
}

const getHeroIconId = (id) => {
    if (!id) return id
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.iconId : id
}

const getHeroName = (id) => {
    if (!id) return ''
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.name : `未知(${id})`
}

const getHeroCountry = (id) => {
    if (!id) return ''
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.country : ''
}

const getHeroType = (id) => {
    if (!id) return ''
    const hero = heroMap[String(resolveHeroId(id))]
    return hero ? hero.type : ''
}

const getSkillName = (id) => {
    if (!id) return ''
    const skill = skillMap[String(id)]
    return skill ? skill.name : `未知(${id})`
}

const getSkillQuality = (id) => {
    if (!id) return ''
    const skill = skillMap[String(id)]
    return skill ? skill.zfQuality : ''
}

const getSkillType = (id) => {
    if (!id) return ''
    const skill = skillMap[String(id)]
    return skill ? skill.type : ''
}

const getTroopTypeId = (team, slot) => {
    if (!team.hero_type) return ''
    const parts = String(team.hero_type).split(',').filter(s => s.trim() !== '')
    // 进攻方过滤第一个，防守方过滤最后一个
    let filtered = team.role === 'attack' ? parts.slice(1) : parts.slice(0, -1)
    // 防守方颠倒顺序
    if (team.role !== 'attack') filtered = filtered.reverse()
    return filtered[slot - 1] ? filtered[slot - 1].trim() : ''
}

const parseSkillInfo = (str, team) => {
    if (!str) return []
    // 格式: 1,id,lv,id,lv,id,lv;2,id,lv,id,lv,id,lv;...
    const groups = String(str).split(';').filter(s => s.trim() !== '')
    const parsed = groups.map(g => {
        const parts = g.split(',')
        return {
            index: parseInt(parts[0]),
            skills: [
                { id: parts[1], level: parseInt(parts[2]) },
                { id: parts[3], level: parseInt(parts[4]) },
                { id: parts[5], level: parseInt(parts[6]) },
            ]
        }
    })
    // 进攻方取 index 1,2,3；防守方取 4,5,6
    let filtered = team.role === 'attack'
        ? parsed.filter(g => g.index >= 1 && g.index <= 3)
        : parsed.filter(g => g.index >= 4 && g.index <= 6)
    // 防守方颠倒
    if (team.role !== 'attack') filtered.reverse()
    return filtered
}

const parseGearInfo = (str, team) => {
    if (!str) return []
    const groups = String(str).split(';').filter(s => s.trim() !== '')
    const parsed = groups.map(g => {
        const parts = g.split(',')
        return {
            gearId: parts[0],
            level: parseInt(parts[1]),
            entryId: parts[2],
        }
    }).filter(g => g.gearId && g.gearId !== '0')
    let filtered = team.role === 'attack' ? parsed : parsed.reverse()
    return filtered
}

const getGearName = (gearId) => {
    if (!gearId || gearId === '0') return ''
    const gear = gearMap[String(gearId)]
    return gear ? gear.name : `未知(${gearId})`
}

const getGearEntryName = (entryId) => {
    if (!entryId || entryId === '0') return ''
    const entry = gearFeatureMap[String(entryId)]
    return entry ? entry.name : `未知(${entryId})`
}

const getGearEntryQuality = (entryId) => {
    if (!entryId || entryId === '0') return ''
    const entry = gearFeatureMap[String(entryId)]
    return entry ? entry.quality : ''
}

const getGearEntryAdvance = (entryId) => {
    if (!entryId || entryId === '0') return 0
    const entry = gearFeatureMap[String(entryId)]
    return entry ? entry.advance : 0
}

const getGearNameClass = (entryId) => {
    const quality = getGearEntryQuality(entryId)
    const advance = getGearEntryAdvance(entryId)
    if (advance === 1) return 'gear-name-red'
    if (quality >= 8) return 'gear-name-pink'
    return 'gear-name-blue'
}

const formatTime = (ts) => {
    if (!ts) return ''
    const d = new Date(ts * 1000)
    const pad = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const groupedResults = () => {
    const map = {}
    results.value.forEach(r => {
        if (!map[r.player_name]) {
            map[r.player_name] = []
        }
        map[r.player_name].push(r)
    })
    return map
}

const roleLabel = (role) => role === 'attack' ? '进攻' : '防守'
const roleType = (role) => role === 'attack' ? 'error' : 'info'
</script>

<template>
    <div class="page-team-query">
        <n-card class="page-card" embedded>
            <div class="page-header">
                <div class="page-header-info">
                    <h2 class="page-title">队伍查询</h2>
                    <p class="page-desc">通过战报数据查询玩家队伍配置</p>
                </div>
            </div>

            <div class="search-bar">
                <n-input v-model:value="searchName" placeholder="玩家名称" clearable @keyup.enter="doSearch" />
                <n-input v-model:value="searchUnion" placeholder="同盟名称" clearable @keyup.enter="doSearch" />
                <n-input v-model:value="searchIdu" placeholder="队伍 ID" clearable @keyup.enter="doSearch" />
                <n-button type="primary" @click="doSearch" :loading="loading">
                    <template #icon><Search :size="16" /></template>
                    查询
                </n-button>
            </div>

            <div class="result-area" v-if="loading">
                <div class="loading-wrap">
                    <n-spin size="medium" />
                    <span>查询中...</span>
                </div>
            </div>

            <div class="result-area" v-else-if="hasSearched && results.length === 0">
                <n-empty description="未找到队伍数据" style="padding: 60px 0;" />
            </div>

            <div class="result-area" v-else-if="results.length > 0">
                <div class="result-summary">
                    共找到 <strong>{{ Object.keys(groupedResults()).length }}</strong> 名玩家，
                    <strong>{{ results.length }}</strong> 支队伍
                </div>

                <div class="player-section" v-for="(teams, playerName) in groupedResults()" :key="playerName">
                    <div class="player-name">
                        <Swords :size="16" />
                        {{ playerName }}
                    </div>

                    <div class="team-card" v-for="team in teams" :key="team.battle_id + team.role + team.hero1_id">
                        <div class="team-header">
                            <n-tag :type="roleType(team.role)" :bordered="false" size="small">
                                {{ roleLabel(team.role) }}
                            </n-tag>
                            <span class="team-idu">
                                {{ team.player_name }} · ID {{ team.idu }}
                            </span>
                            <span class="team-stars">
                                <Star :size="13" />
                                红度 {{ team.total_star }}
                            </span>
                            <span class="team-time">{{ formatTime(team.time) }}</span>
                        </div>

                        <div class="hero-row">
                            <div class="hero-slot" v-for="i in 3" :key="i">
                                <div class="hero-avatar">
                                    <img v-if="team[`hero${i}_id`]"
                                        :src="`https://g0.gph.netease.com/ngsocial/community/stzb/cn/cards/cut/card_small_${getHeroIconId(team[`hero${i}_id`])}.jpg?gameid=g10`"
                                        @error="$event.target.style.display='none'" />
                                    <div class="hero-placeholder" v-else>?</div>
                                </div>
                                <div class="hero-info">
                                    <span class="hero-name">{{ getHeroName(team[`hero${i}_id`]) }}</span>
                                    <span class="hero-meta">
                                        <n-tag v-if="getHeroCountry(team[`hero${i}_id`])" size="tiny" :bordered="false">{{ getHeroCountry(team[`hero${i}_id`]) }}</n-tag>
                                        <n-tag v-if="getHeroType(team[`hero${i}_id`])" size="tiny" :bordered="false" type="info">{{ getHeroType(team[`hero${i}_id`]) }}</n-tag>
                                        <span class="hero-level">Lv.{{ team[`hero${i}_level`] }}</span>
                                        <span class="hero-star">{{ team[`hero${i}_star`] }}红</span>
                                    </span>
                                </div>
                                <div class="troop-type-wrap" v-if="getTroopTypeId(team, i)">
                                    <img class="troop-type-img" :src="`https://cbg-stzb.res.netease.com/mvvm/rc346663d4140700aaab6da137/images/bz/${getTroopTypeId(team, i)}.png`" @error="$event.target.style.display='none'" />
                                </div>
                            </div>
                        </div>

                        <div class="skill-section" v-if="team.all_skill_info">
                            <div class="skill-hero-row">
                                <div class="skill-hero-card" v-for="(group, gi) in parseSkillInfo(team.all_skill_info, team)" :key="gi">
                                    <template v-for="(skill, si) in group.skills" :key="si">
                                    <div class="skill-slot" v-if="skill && skill.id && skill.id !== '0'">
                                        <span class="skill-meta">
                                            <n-tag v-if="getSkillQuality(skill.id)" size="tiny" :bordered="false" :type="getSkillQuality(skill.id) === 'S' ? 'warning' : getSkillQuality(skill.id) === 'A' ? 'info' : 'default'">{{ getSkillQuality(skill.id) }}</n-tag>
                                            <n-tag v-if="getSkillType(skill.id)" size="tiny" :bordered="false">{{ getSkillType(skill.id) }}</n-tag>
                                            <span class="skill-name">{{ getSkillName(skill.id) }}</span>
                                            <span class="skill-level">Lv.{{ skill.level }}</span>
                                        </span>
                                    </div>
                                    </template>
                                </div>
                            </div>
                        </div>

                        <div class="gear-section" v-if="team.gear">
                            <div class="gear-row">
                                <div class="gear-card" v-for="(gear, gi) in parseGearInfo(team.gear, team)" :key="gi">
                                    <div class="gear-img-wrap">
                                        <img class="gear-img" :src="`https://cbg-stzb.res.netease.com/game_res/gears/gear_icon/gear_icon_${gear.gearId}.jpg`" @error="$event.target.style.display='none'" />
                                    </div>
                                    <div class="gear-info">
                                        <span class="gear-meta">
                                            <span class="gear-base-name">{{ getGearName(gear.gearId) }}</span>
                                            <span class="gear-name" :class="getGearNameClass(gear.entryId)">{{ getGearEntryName(gear.entryId) }}</span>
                                            <span class="gear-level">Lv.{{ gear.level }}</span>
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </n-card>
    </div>
</template>

<style scoped lang="scss">
.page-team-query {
    display: flex;
    flex-direction: column;
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

.search-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
}

.search-bar .n-input {
    flex: 1;
    min-width: 160px;
}

.result-area {
    min-height: 200px;
}

.loading-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 60px 0;
    color: #64748b;
    font-size: 14px;
}

.result-summary {
    font-size: 13px;
    color: #64748b;
    margin-bottom: 16px;
}

.player-section {
    margin-bottom: 24px;
}

.player-name {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 700;
    color: #1e293b;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 2px solid #e2e8f0;
}

.team-card {
    background: #fff;
    border: 1px solid rgba(228, 228, 231, 0.6);
    border-radius: 10px;
    padding: 16px;
    margin-bottom: 12px;
    transition: box-shadow 0.2s;

    &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
    }
}

.team-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
    flex-wrap: wrap;
}

.team-idu,
.team-stars,
.team-time {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: #64748b;
}

.team-time {
    margin-left: auto;
    font-size: 12px;
}

.hero-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin-bottom: 12px;
}

.hero-slot {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    background: #f8f9fb;
    border-radius: 8px;
}

.troop-type-wrap {
    margin-left: auto;
    width: 36px;
    height: 36px;
    border-radius: 8px;
    background: #1e293b;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
}

.troop-type-img {
    width: 28px;
    height: 28px;
    object-fit: contain;
}

.hero-avatar {
    width: 44px;
    height: 44px;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
    background: #e2e8f0;

    img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }
}

.hero-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    color: #94a3b8;
}

.hero-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.hero-name {
    font-size: 14px;
    font-weight: 600;
    color: #1e293b;
}

.hero-meta {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-wrap: wrap;
}

.hero-level {
    font-size: 12px;
    font-weight: 600;
    color: #1e293b;
}

.hero-star {
    font-size: 12px;
    color: #f59e0b;
    font-weight: 500;
}

.skill-section {
    margin-bottom: 12px;
}

.skill-hero-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
}

.skill-hero-card {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 10px;
    background: #f8f9fb;
    border-radius: 8px;
}

.skill-slot {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.skill-name {
    font-size: 13px;
    font-weight: 600;
    color: #1e293b;
}

.skill-meta {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-wrap: wrap;
}

.skill-level {
    font-size: 12px;
    color: #64748b;
    font-weight: 500;
}

.skill-detail {
    max-height: 300px;
    overflow-y: auto;
}

.json-block {
    font-family: 'Cascadia Code', 'Fira Code', 'Consolas', monospace;
    font-size: 12px;
    line-height: 1.5;
    color: #334155;
    margin: 0;
    white-space: pre-wrap;
    word-break: break-all;
    background: #f8f9fb;
    padding: 12px;
    border-radius: 6px;
}

.gear-section {
    margin-top: 4px;
}

.gear-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
}

.gear-card {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    background: #f8f9fb;
    border-radius: 8px;
}

.gear-img-wrap {
    width: 44px;
    height: 44px;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
    background: #e2e8f0;

    img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }
}

.gear-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.gear-base-name {
    font-size: 14px;
    font-weight: 600;
    color: #1e293b;
}

.gear-name {
    font-size: 12px;
    font-weight: 500;
    color: #fff;
    padding: 0 4px;
    border-radius: 3px;
    display: inline-block;
}

.gear-name-blue {
    background: #3b82f6;
}

.gear-name-pink {
    background: #ec4899;
}

.gear-name-red {
    background: #ef4444;
}

.gear-meta {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-wrap: wrap;
}

.gear-level {
    font-size: 12px;
    color: #64748b;
    font-weight: 500;
}
</style>
