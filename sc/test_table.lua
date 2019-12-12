function run(tbl)
    if tbl["requesttime"] == -1 then
        if tbl["gender"] == "man" then
            return "你好，先⽣"
        elseif tbl["gender"] == "woman" then
            return "你好，女士"
        else
            return "你好"
        end
    else
        if tbl["requesttime"] >= 600 and tbl["requesttime"] <= 1200 then
            if tbl["gender"] == "man" then
                return "早上好，先⽣"
            elseif tbl["gender"] == "woman" then
                return "早上好，女士"
            else
                return "早上好"
            end
        elseif tbl["requesttime"] > 1200 and tbl["requesttime"] <= 1400 then
            if tbl["gender"] == "man" then
                return "中午好，先⽣"
            elseif tbl["gender"] == "woman" then
                return "中午好，女士"
            else
                return "中午好"
            end
        elseif tbl["requesttime"] > 1400 and tbl["requesttime"] < 1800 then
            if tbl["gender"] == "man" then
                return "下午好，先⽣"
            elseif tbl["gender"] == "woman" then
                return "下午好，女士"
            else
                return "下午好"
            end

        else
            if tbl["gender"] == "man" then
                return "晚上好，先⽣"
            elseif tbl["gender"] == "woman" then
                return "晚上好，女士"
            else
                return "晚上好"
            end
        end
    end
end

