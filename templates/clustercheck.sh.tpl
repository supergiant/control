until $([ $(sudo kubectl get nodes|grep Ready|grep master|wc -l) -eq 1 ]); do printf '.'; sleep 5; done

until $([ $(sudo kubectl get nodes|grep Ready|grep none|wc -l) -eq 1 ]); do printf '.'; sleep 5; done
