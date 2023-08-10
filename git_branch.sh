#!/bin/bash
echo "<<<<<<<< Git自动化布署开始 >>>>>>>>"


function check_result() {
  makeRet=$1
  message=$2
  if [ "$makeRet" -ne 0 ]; then
    echo "====>>>>" "${message}"
    exit "$makeRet" # end the script running
  fi
}

# exit
execExit() {
  echo "<<<<<<<< 中断程序执行 >>>>>>>>"
  exit 1
}

# 新建分支
checkoutNewBranch() {
  newBranch=$1
  if [ "$newBranch" != "" ];then
    git checkout -b ${newBranch}
  else
    echo "输入值为空"
    execExit
  fi
}

# 切换分支
checkoutBranch() {
  # 输入提交说明
  branchName=$1
  git checkout ${branchName}
}

# 拉取远程仓库分支代码
pullRepo() {
  git add .
  git status

  # 输入提交说明
  read -p "请输入本次提交的备注说明:" commit
  echo "<<<<<<<< 将暂存区内容提交到本地仓库:开始 >>>>>>>>"
  git commit -m ${commit} --no-verify

  # 输入拉取远程仓库的分支名称
  read -p "请输入要推送到远程仓库的分支名称:" pullBranch
  if [ "$pullBranch" != "" ]; then
    echo "<<<<<<<< 拉取远程分支到本地仓库并合并:${pullBranch}开始 >>>>>>>>"
    git pull origin ${pullBranch}
  else
    git pull
  fi
}

# 推送到远程分支代码
pushRepo() {
  # 输入推到远程仓库的分支名称
  pushBranch=$1
  if [ "$pushBranch" != "" ]; then
    echo "<<<<<<<< 推送本地分支更到远程分支并合并:${pullBranch}开始 >>>>>>>>"
    git push origin ${pushBranch}
  else
    git push
  fi
}

# 当前分支
main(){
  newBranch=$1
  old_branch=`git  branch | grep '*' | sed -e 's/\*//g' -e 's/HEAD detached at//g' -e 's/\s*//g' -e 's/[\(\)]//g'`

  echo "开始新建分支:${newBranch}"
  checkoutNewBranch  "${newBranch}"
  check_result $? "checkout NewBranch err"
  echo "提交新分支:${newBranch}到远程"
  pushRepo "${newBranch}"
  check_result $? "NewBranch  pushRepo err"
  echo "切换到原来的分支:${old_branch}"
  checkoutBranch "${old_branch}"
  check_result $? "checkout old Branch err"
  echo "打分支完成"
}


main $1




