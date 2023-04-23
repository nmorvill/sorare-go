import data from "./export.json" assert { type: 'json' };

export function getCalendarsTable() {
    return createRows(data)
}

function createRows(teams) {
    let ret = ""
    for(const team of teams) {
        ret += createRow(team)
    }
    return ret
}

function createRow(team) {
    return `
    <div class="row">
        <div class="row-head">
            <img src="${team.logoURL}"/>
            <div class="names">
                <h3>${team.abbr}</h3>
                ${team.name}
            </div>
        </div>
        ${createCells(team.games, team.nbTeams)}
    </div>
    `
}

function createCells(games, nbTeams) {
    let ret = ""
    for(let i = 0; i < 5; i++) {
        if(games != null && i < games.length ) {
            ret += createCell(games[i], nbTeams)
        } else {
            ret += createEmptyCell()
        }
    }
    return ret
}

function createCell(game, nbTeams) {
    return `
    <div class="cell" style="background-color:${getBackgroundColor(game.oppRank, nbTeams)}">
        <img src="${game.logoURL}"/>
        ${game.oppRank}
        <div class="location ${game.location == "HOME" ? "home" : "away"}">
            ${game.location == "HOME" ? "H" : "A"}
        </div>
    </div>
    `
}

function createEmptyCell() {
    return `
    <div class="cell" style="background-color:#000000">
    </div>
    `
}

function getBackgroundColor(rank, maxRank) {
    const green = {
        r : 87,
        g : 223,
        b: 72
    }
    const red = {
        r : 255,
        g : 67,
        b : 57
    }
    const per = rank/maxRank
    const r = red.r + per * (green.r - red.r)
    const g = red.g + per * (green.g - red.g)
    const b = red.b + per * (green.b - red.b)
    return `rgb(${r}, ${g}, ${b}`
}
