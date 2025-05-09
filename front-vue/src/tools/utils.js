//路径拼接
//('/root', 'images', 'logo.png')  // → "/root/images/logo.png"
//('/root', '', 'logo.png')        // → "/root/logo.png"



export const buildPath=(base,sub,filename)=>{
    //
    return sub ? `${base}/${sub}/${filename}` : `${base}/${filename}`;
}


//buildFolderPath('/data', 'docs', '2024')  // → "/data/docs/2024"
//buildFolderPath('/data', '', '2024')      // → "/data/2024"

export const buildFolderPath = (base, sub, folder) => {
    return sub ? `${base}/${sub}/${folder}` : `${base}/${folder}`;
  };
  